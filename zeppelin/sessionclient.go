package zeppelin

import (
	"fmt"
	"strings"

	"github.com/DataWorkbench/common/qerror"
)

type ZSession struct {
	zeppelinClient *Client
	interpreter    string
	intpProperties map[string]string
	maxStatement   int
	sessionInfo    *SessionInfo
}

func NewZSession(config ClientConfig, interceptor string) *ZSession {
	return NewZSessionWithProperties(config, interceptor, make(map[string]string))
}

func NewZSessionWithSessionId(config ClientConfig, interceptor string, sessionId string) (*ZSession, error) {
	sessionInfo, err := NewSessionInfo([]byte(sessionId))
	if err != nil {
		return nil, err
	}
	return &ZSession{
		zeppelinClient: NewZeppelinClient(config),
		interpreter:    interceptor,
		sessionInfo:    sessionInfo,
	}, nil
}

func NewZSessionWithProperties(config ClientConfig, interceptor string, intpPorperties map[string]string) *ZSession {
	return NewZSessionWithAll(config, interceptor, intpPorperties, 100)
}

func NewZSessionWithAll(config ClientConfig, interceptor string, intpPorperties map[string]string, maxStatement int) *ZSession {
	return &ZSession{
		zeppelinClient: NewZeppelinClient(config),
		interpreter:    interceptor,
		intpProperties: intpPorperties,
		maxStatement:   maxStatement,
	}
}

func CreateFromExistingSession(config ClientConfig, interceptor string, sessionId string) (*ZSession, error) {
	session, err := NewZSessionWithSessionId(config, interceptor, sessionId)
	if err != nil {
		return nil, err
	}
	if err = session.Reconnect(); err != nil {
		return nil, err
	}
	return session, nil
}

func (z *ZSession) Start() (err error) {
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
	if !paragraphResult.Status.IsFinished() {
		return qerror.ZeppelinConfigureFailed
	}

	if paragraphId, err = z.zeppelinClient.addParagraph(z.getNoteId(), "Session Init", "%"+z.interpreter+"(init=true)"); err != nil {
		return
	}
	if paragraphResult, err = z.zeppelinClient.executeParagraphWithSessionId(z.getNoteId(), paragraphId, z.getSessionId()); err != nil {
		return
	}
	if !paragraphResult.Status.IsFinished() {
		return qerror.ZeppelinInitFailed
	}
	return nil
}

func (z *ZSession) Stop() error {
	if z.getSessionId() != "" {
		return z.zeppelinClient.stopSession(z.getSessionId())
	}
	return nil
}

func (z *ZSession) SubmitWithProperties(subInterpreter string, localProperties map[string]string, code string) (*ExecuteResult, error) {
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

func (z *ZSession) Submit(subInterpreter string, code string) (*ExecuteResult, error) {
	return z.SubmitWithProperties(subInterpreter, make(map[string]string), code)
}

func (z *ZSession) Sub(code string) (*ExecuteResult, error) {
	return z.Submit("", code)
}

func (z *ZSession) ExecuteWithProperties(subInterpreter string, localProperties map[string]string, code string) (*ExecuteResult, error) {
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

func (z *ZSession) Execute(subInterpreter string, code string) (*ExecuteResult, error) {
	return z.ExecuteWithProperties(subInterpreter, make(map[string]string), code)
}

func (z *ZSession) Exec(code string) (*ExecuteResult, error) {
	return z.Execute("", code)
}

func (z *ZSession) Cancel(statementId string) error {
	return z.zeppelinClient.cancelParagraph(z.getNoteId(), statementId)
}

func (z *ZSession) QueryStatement(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.queryParagraphResult(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) WaitUntilFinished(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.waitUtilParagraphFinish(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) WaitUntilRunning(statementId string) (*ExecuteResult, error) {
	paragraphResult, err := z.zeppelinClient.waitUtilParagraphRunning(z.getNoteId(), statementId)
	if err != nil {
		return nil, err
	}
	return NewExecuteResult(paragraphResult), nil
}

func (z *ZSession) Reconnect() (err error) {
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
