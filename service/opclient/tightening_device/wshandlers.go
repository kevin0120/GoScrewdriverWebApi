package tightening_device

import (
	"encoding/json"
	"github.com/masami10/rush/utils"
	"time"

	"github.com/kataras/iris/v12/websocket"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/wsnotify"
	uuid "github.com/satori/go.uuid"
)

var getPSetInfoWithDetail = utils.GetEnvBool("ENV_PSET_INFO_DETAIL", true)

func (s *Service) OnWS_TOOL_MODE_SELECT(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolModeSelect{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolModeSelect(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_ENABLE(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	req := ToolControl{}
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolControl(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

	//FIXME: 此功能未来使用nodered控制 使能/切断使能后，控制工具对应LED灯。
	s.doDispatch(dispatcherbus.DispatcherToolEnable, req)

	if !req.Enable && s.config().SocketSelector.Enable {
		s.doDispatch(dispatcherbus.DispatcherSocketSelector, SocketSelectorReq{
			Type: SocketSelectorClear,
		})
	}
}

func (s *Service) OnWS_TOOL_JOB(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req JobSet
	_ = json.Unmarshal(byteData, &req)
	err := s.ToolJobSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

}

func (s *Service) OnWS_TOOL_PSET(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetSet
	if err := json.Unmarshal(byteData, &req); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
	}

	if s.config().SocketSelector.Enable {
		//如果使能套筒选择器，不设定程序号，而是控制套筒选择器
		s.doDispatch(dispatcherbus.DispatcherSocketSelector, SocketSelectorReq{
			PSetSet: req,
			Type:    SocketSelectorTrigger,
		})
	} else {
		err := s.ToolPSetSet(&req)
		if err != nil {
			_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
			return
		}
	}

	//if req.Enable || ENV_PSET_WITH_ENABLE {
	//	req.ToolControl.Enable = true
	//	if err := s.ToolControl(&req.ToolControl); err != nil {
	//		s.diag.Error("OnWS_TOOL_PSET ToolControl", err)
	//		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
	//	}
	//}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)

}

func (s *Service) OnWS_TOOL_PSET_BATCH(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req PSetBatchSet
	_ = json.Unmarshal(byteData, &req)

	err := s.ToolPSetBatchSet(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_LIST(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)
	var pSetList interface{}
	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	pSetNums, err := s.GetToolPSetList(&req)

	//TODO 砺星扳手暂时不支持获取详情
	controller, _ := s.getController(req.ControllerSN)
	if getPSetInfoWithDetail && (controller.Model() != ModelLexenWrench) {
		var _pSetList []*PSetDetail
		for i := 0; i < len(pSetNums); i++ {
			pSetInfo := pSetNums[i]
			pSetDetail, err1 := s.GetToolPSetDetail(&ToolPSet{
				req,
				pSetInfo.ID,
			})
			if err1 != nil {
				s.diag.Error("error when GetToolPSetDetail", err1)
				continue
			}
			_pSetList = append(_pSetList, pSetDetail)
		}
		pSetList = _pSetList
	} else {
		var _pSetList []*PSetDetail
		for i := 0; i < len(pSetNums); i++ {
			_pSetList = append(_pSetList, &PSetDetail{
				PSetID:   pSetNums[i].ID,
				PSetName: pSetNums[i].Name,
			})
		}
		pSetList = _pSetList
	}

	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, pSetList), s.diag)
}

func (s *Service) OnWS_TOOL_PSET_DETAIL(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolPSet
	_ = json.Unmarshal(byteData, &req)

	psetDetail, err := s.GetToolPSetDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, psetDetail), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_LIST(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolInfo
	_ = json.Unmarshal(byteData, &req)

	jobList, err := s.GetToolJobList(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, jobList), s.diag)
}

func (s *Service) OnWS_TOOL_JOB_DETAIL(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var req ToolJob
	_ = json.Unmarshal(byteData, &req)

	jobDetail, err := s.GetToolJobDetail(&req)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateWSMsg(msg.SeqNumber, msg.Type, jobDetail), s.diag)
}

//手动回补填写结果
func (s *Service) OnWS_TOOL_RESULT_MANUAL_SET(c *websocket.NSConn, msg *wsnotify.WSMsg) {
	byteData, _ := json.Marshal(msg.Data)

	var result TighteningResult
	_ = json.Unmarshal(byteData, &result)

	if err := result.ValidateSet(); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -1, err.Error()), s.diag)
		return
	}

	tool, err := s.getTool(result.ControllerSN, result.ToolSN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -2, err.Error()), s.diag)
		return
	}

	dbTool, err := s.storageService.GetTool(result.ToolSN)
	if err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -3, err.Error()), s.diag)
		return
	}

	dbTool.Count = result.Count
	//fixme: xorm 无法设置0的情况
	if result.Count <= 0 {
		dbTool.Count = 1
	}
	if err := s.storageService.UpdateTool(&dbTool); err != nil {
		_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, -4, err.Error()), s.diag)
		return
	}
	// 未设置userid会导致数据插入失败
	if result.UserID == 0 {
		result.UserID = 1
	}
	result.TighteningUnit = result.ToolSN
	result.TighteningID = uuid.NewV4().String()
	//手动输入时间以rush收到为准
	result.UpdateTime = time.Now()
	// 处理数据
	result.StepResults = []StepData{
		StepData{
			PSetDefine{}, result,
		},
	}
	s.doDispatch(tool.GenerateDispatcherNameBySerialNumber(dispatcherbus.DispatcherResult), &result)
	_ = wsnotify.WSClientSend(c, wsnotify.WS_EVENT_REPLY, wsnotify.GenerateReply(msg.SeqNumber, msg.Type, 0, ""), s.diag)
}
