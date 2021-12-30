package zeppelin

import (
	"fmt"
	"strings"

	"github.com/DataWorkbench/common/qerror"
)

type ZSession struct {
	zeppelinClient  *Client
	interpreter     string
	intpProperties  map[string]string
	maxStatement    int
	webSocketClient *WebSocketClient
	sessionInfo     *SessionInfo
}

func NewZSession2(config ClientConfig, interceptor string) *ZSession {
	return NewZSession4(config, interceptor, make(map[string]string), 100)
}

func NewZSession3(config ClientConfig, interceptor string, sessionId string) *ZSession {
	sessionInfo, _ := NewSessionInfo(sessionId)
	return &ZSession{
		zeppelinClient: NewZeppelinClient(config),
		interpreter:    interceptor,
		sessionInfo:    sessionInfo,
	}
}

func NewZSession4(config ClientConfig, interceptor string, intpPorperties map[string]string, maxStatement int) *ZSession {
	return &ZSession{
		zeppelinClient: NewZeppelinClient(config),
		interpreter:    interceptor,
		intpProperties: intpPorperties,
		maxStatement:   maxStatement,
	}
}

func CreateFromExistingSession(config ClientConfig, interceptor string, sessionId string) (*ZSession, error) {
	session := NewZSession3(config, interceptor, sessionId)
	if err := session.reconnect(); err != nil {
		return nil, err
	}
	return session, nil
}

func (z *ZSession) start() (err error) {
	if z.sessionInfo, err = z.zeppelinClient.newSession(z.interpreter); err != nil {
		return
	}
	var builder strings.Builder
	builder.WriteString("%" + z.interpreter + ".conf\n")
	if z.intpProperties != nil {
		for k, v := range z.intpProperties {
			builder.WriteString(k + " " + v + "\n")
		}
	}
	var (
		paragraphId     string
		paragraphResult *ParagraphResult
	)
	if paragraphId, err = z.zeppelinClient.addParagraph(z.getNoteId(), "Session Configuration", builder.String()); err != nil {
		return
	}
	if paragraphResult, err = z.zeppelinClient.executeParagraphWithSessionId(z.getNoteId(), paragraphId, z.getSessionId()); err != nil {
		return
	}
	if !paragraphResult.Status.isFinished() {
		return qerror.ZeppelinConfigureFailed
	}

	if paragraphId, err = z.zeppelinClient.addParagraph(z.getNoteId(), "Session Init", "%"+z.interpreter+"(init=true)"); err != nil {
		return
	}
	if paragraphResult, err = z.zeppelinClient.executeParagraphWithSessionId(z.getNoteId(), paragraphId, z.getSessionId()); err != nil {
		return
	}
	if !paragraphResult.Status.isFinished() {
		return qerror.ZeppelinInitFailed
	}
	/*if z.sessionInfo, err = z.zeppelinClient.getSession(z.getSessionId()); err != nil {
		return
	}
	if handler != nil {
		z.webSocketClient = NewWebSocketClient(handler)
		restUrl := z.zeppelinClient.ClientConfig.ZeppelinRestUrl
		wsUrl := strings.ReplaceAll(restUrl, "https", "ws")
		wsUrl = strings.ReplaceAll(wsUrl, "http", "ws") + "/ws"
		req := map[string]string{}
		req["id"] = z.getNoteId()
		req["op"] = "GET_NOTE"
		var reqBytes []byte
		if reqBytes, err = json.Marshal(req); err != nil {
			return
		}
		if err = z.webSocketClient.connect(wsUrl); err != nil {
			return
		}
		return z.webSocketClient.dial.WriteJSON(string(reqBytes))
	}*/
	return nil
}

func (z *ZSession) stop() error {
	if z.getSessionId() != "" {
		return z.zeppelinClient.stopSession(z.getSessionId())
	}
	//if z.webSocketClient != nil {
	//	//TODO stop websocket
	//}
	return nil
}

func (z *ZSession) submitWithProperties(subInterpreter string, localProperties map[string]string, code string) (*ExecuteResult, error) {
	builder := strings.Builder{}
	builder.WriteString("%" + z.interpreter)
	if subInterpreter != "" && len(subInterpreter) > 0 {
		builder.WriteString("." + subInterpreter)
	}
	if localProperties != nil && len(localProperties) > 0 {
		builder.WriteString("(")
		var propertyStrs []string
		for k, v := range localProperties {
			propertyStrs = append(propertyStrs, fmt.Sprintf("\"%s\"=\"%s\"", k, v))
		}
		builder.WriteString(strings.Join(propertyStrs, ","))
		builder.WriteString(")")
	}
	builder.WriteString(" " + code)
	text := builder.String()
	nextParagraphId, err := z.zeppelinClient.addParagraph(z.getNoteId(), "", text)
	if err != nil {
		return nil, err
	}
	paragraphResult, err := z.zeppelinClient.submitParagraphWithSessionId(z.getNoteId(), nextParagraphId, z.getSessionId())
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) submit(subInterpreter string, code string) (*ExecuteResult, error) {
	return z.submitWithProperties(subInterpreter, make(map[string]string), code)
}

func (z *ZSession) executeWithProperties(subInterpreter string, localProperties map[string]string, code string) (*ExecuteResult, error) {
	builder := strings.Builder{}
	builder.WriteString("%" + z.interpreter)
	if subInterpreter != "" && len(subInterpreter) > 0 {
		builder.WriteString("." + subInterpreter)
	}
	if localProperties != nil && len(localProperties) > 0 {
		builder.WriteString("(")
		var propertyStrs []string
		for k, v := range localProperties {
			propertyStrs = append(propertyStrs, fmt.Sprintf("\"%s\"=\"%s\"", k, v))
		}
		builder.WriteString(strings.Join(propertyStrs, ","))
		builder.WriteString(")")
	}
	builder.WriteString(" " + code)
	text := builder.String()
	nextParagraphId, err := z.zeppelinClient.addParagraph(z.getNoteId(), "", text)
	if err != nil {
		return nil, err
	}
	paragraphResult, err := z.zeppelinClient.executeParagraphWithSessionId(z.getNoteId(), nextParagraphId, z.getSessionId())
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) cancel(statementId string) error {
	return z.zeppelinClient.cancelParagraph(z.getNoteId(), statementId)
}

func (z *ZSession) queryStatement(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.queryParagraphResult(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) waitUntilFinished(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.waitUtilParagraphFinish(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) waitUntilRunning(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.waitUtilParagraphRunning(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) reconnect() (err error) {
	z.sessionInfo, err = z.zeppelinClient.getSession(z.getSessionId())
	if !strings.EqualFold("Running", z.sessionInfo.State) {
		return qerror.ZeppelinSessionNotRunning
	}
	return nil
}

func (z *ZSession) getNoteId() string {
	if z.sessionInfo != nil {
		return z.sessionInfo.NoteId
	}
	return ""
}

func (z *ZSession) getSessionId() string {
	if z.sessionInfo != nil {
		return z.sessionInfo.SessionId
	}
	return ""
}
