// Package ai_sdk
// @Author Clover
// @Data 2024/8/11 下午10:34:00
// @Desc 多会话管理器
package ai_sdk

import (
	"fmt"
	"github.com/Clov614/go-ai-sdk/config"
	"sync"
	"time"
)

const (
	systemRole    = "system"
	userRole      = "user"
	assistantRole = "assistant"
)

// Session 会话主体 （k-v 会话id-会话信息）
type Session struct {
	globalSurvivalLimit time.Duration
	systemContent       string // 预设消息
	cache               map[string]*sessionInfo
	mu                  sync.Mutex
}

func NewSession(systemSet string, persessionTimeOut int) *Session {
	sessionTimeOut := time.Duration(persessionTimeOut) * time.Minute
	if sessionTimeOut < 2*time.Minute {
		sessionTimeOut = 2 * time.Minute
	}
	s := Session{
		globalSurvivalLimit: sessionTimeOut,
		systemContent:       systemSet,
		cache:               make(map[string]*sessionInfo, 10),
	}
	return &s
}

// 会话信息
type sessionInfo struct {
	sessionId      string        // 会话唯一id
	history        *history      // history: 上下文
	startTime      time.Time     // 会话时间信息
	survivalLimit  time.Duration // 会话存活时间(存活周期)
	survivalSignal chan struct{} // 存活信号量
	done           chan struct{} // 结束信号
}

// GetSession 获取唯一会话
func (s *Session) GetSession(sessionId string) *sessionInfo {
	s.mu.Lock()
	if info, ok := s.cache[sessionId]; ok { // 判断会话列表中是否存在该id（存在则不创建，直接返回）
		s.mu.Unlock()
		return info
	}
	s.mu.Unlock()
	return s.newSession(sessionId)
}

// 新创建的会话启动计时器，并通过存活信号量刷新计时器，超时移除该会话(注意goroutine泄露问题)
func (s *Session) newSession(sessionId string) *sessionInfo {
	s.mu.Lock()
	info := &sessionInfo{
		sessionId:      sessionId,
		history:        newHistory(s.systemContent),
		startTime:      time.Now(),
		survivalLimit:  s.globalSurvivalLimit,
		survivalSignal: make(chan struct{}),
		done:           make(chan struct{}),
	}
	s.cache[sessionId] = info
	s.mu.Unlock()
	s.checkSurvival(info) // 超时检测
	return info
}

// 超时检测
func (s *Session) checkSurvival(info *sessionInfo) {
	var once sync.Once
	go func() {
		timer := time.NewTimer(info.survivalLimit)
		defer timer.Stop()
		for {
			select {
			case <-info.survivalSignal: // 重置计时器
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(info.survivalLimit)
			case <-timer.C: // 超时删除会话
				once.Do(func() {
					s.removeById(info.sessionId)
				})
				return
			case <-info.done:
				once.Do(func() {
					s.removeById(info.sessionId)
				})
				return
			}
		}
	}()
}

// TalkById 根据会话id对话 新增会话来获取会话发起对话
func (s *Session) TalkById(sessionId string, content string) (string, error) {
	sessioninfo := s.GetSession(sessionId)
	return sessioninfo.Talk(content)
}

// Talk 对该sessionInfo 发起对话
func (s *sessionInfo) Talk(content string) (string, error) {
	go func() {
		s.survivalSignal <- struct{}{} // 确保在会话期间存活
	}()
	answer, err := s.history.handleQuestion(content, func(msgs []ChatMessage) (answer ChatMessage, err error) {
		resp, err := aiclient.Send(Request{Messages: msgs})
		if err != nil {
			return answer, fmt.Errorf("aiclient.Send err: %w", err)
		}
		answer.Role = assistantRole
		answer.Content = resp.data.Choices[0].Message.Content
		return answer, nil
	})
	if err != nil {
		return "", fmt.Errorf("talkById err: %w", err)
	}
	return answer.Content, nil
}

// 移除会话
func (s *Session) removeById(id string) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok = s.cache[id]; ok {
		delete(s.cache, id)
	}
	return
}

// history 上下文
type history struct {
	maxHistory int // 最长上下文数量
	system     ChatMessage
	dialog     []dialogEntry // 问答实体类
	msgListNum int           // 当前上下文长度
	msgList    []ChatMessage // 转换后请求用上下文
	mu         sync.Mutex
}

func newHistory(system string) *history {
	h := history{
		maxHistory: config.Config.HistoryNum,
		dialog:     make([]dialogEntry, 0),
		msgListNum: 0,
		msgList:    make([]ChatMessage, 0),
	}
	if system != "" {
		h.system = ChatMessage{
			Role:    systemRole,
			Content: system,
		}
	}
	return &h
}

// 问答实体类
type dialogEntry struct {
	question ChatMessage
	answer   ChatMessage
}

// 处理问题
func (h *history) handleQuestion(content string, handleFunc func(msgs []ChatMessage) (answer ChatMessage, err error)) (answer ChatMessage, err error) {
	msgs := h.getChatMessage()
	question := ChatMessage{
		Role:    userRole,
		Content: content,
	}
	answer, err = handleFunc(append(msgs, question))
	if err != nil {
		return answer, fmt.Errorf("handleQuestion handleMessage err: %w", err)
	}
	h.addLast(dialogEntry{question: question, answer: answer})
	return answer, err
}

// concurrent unsafe 删除最早的一条对话记录
func (h *history) removeFirst() (removedEntry dialogEntry) {
	if len(h.dialog) == 0 {
		return dialogEntry{}
	}

	removeEntry := (h.dialog)[0]
	if cap(h.dialog) > 64 {
		var newDialog []dialogEntry
		copy(h.dialog[1:], newDialog)
		h.dialog = newDialog
		return removeEntry
	}
	h.dialog = h.dialog[1:]
	return removeEntry
}

func (h *history) addLast(entry dialogEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for len(h.dialog) >= h.maxHistory {
		h.removeFirst()
	}
	h.dialog = append(h.dialog, entry)
}

func (h *history) getChatMessage() (msg []ChatMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	dialogLen := len(h.dialog)
	if dialogLen == 0 {
		h.msgList = append(h.msgList, h.system)
	}
	if h.msgListNum == dialogLen { // 长度相同，直接返回
		return h.msgList
	}
	for h.msgListNum < dialogLen {
		// 添加会话
		h.msgList = append(h.msgList, h.dialog[h.msgListNum].question)
		h.msgList = append(h.msgList, h.dialog[h.msgListNum].answer)
		h.msgListNum++
	}
	// 返回一个 h.msgList 的副本
	return append([]ChatMessage(nil), h.msgList...)
}

//func (h *history) removeLastQuestion() (question ChatMessage, ok bool) {
//	h.mu.Lock()
//	defer h.mu.Unlock()
//	last := (*h.msgList)[len(*h.msgList) - 1]
//	if last.Role == user_role {
//		question = last
//		*h.msgList = (*h.msgList)[:len(*h.msgList) - 1] // 删除最后一条问题
//		ok = true
//	}
//	return
//}

var DefaultSession *Session

const defaultSystemSet = "从现在开始，我需要你扮演小鸟游六花这个动漫角色，语气跟说话逻辑都要尽力模仿，要完美融入这个角色的设定中。我会称呼你为rikka或六花，届时你明白是在称呼你。"

func init() {
	sessionTimeOut := time.Duration(config.Config.SessionTimeOut) * time.Minute
	if sessionTimeOut < 2*time.Minute {
		sessionTimeOut = 2 * time.Minute
	}
	DefaultSession = &Session{
		globalSurvivalLimit: sessionTimeOut,
		systemContent:       defaultSystemSet,
		cache:               make(map[string]*sessionInfo, 10),
	}
}
