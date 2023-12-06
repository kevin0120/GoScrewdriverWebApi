package openprotocol

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/masami10/rush/utils/ascii"
	"github.com/stretchr/testify/assert"
)

func Test_parseOpenProtocolErrorCode(t *testing.T) {
	err := "111"
	ss, _ := parseOpenProtocolErrorCode(err)
	t.Log(ss)
}

func procHandleRecv() {
	//c.receiveBuf = make(chan []byte, BufferSize)
	//handleRecvBuf := make([]byte, BufferSize)
	//lenBuf := len(handleRecvBuf)
	handleRecvBuf := []byte{}
	var writeOffset = 0
	buf := []byte{1, 2, 3, 4, 5, 6, 7, OpTerminal, 1, 2, 3, 4, 5, 6, 7, OpTerminal}

	for {
		// 处理接收缓冲
		var readOffset = 0
		var index = 0

		for {
			if readOffset >= len(buf) {
				break
			}
			index = bytes.IndexByte(buf[readOffset:], OpTerminal)
			if index == -1 {
				// 没有结束字符,放入缓冲等待后续处理
				restBuf := buf[readOffset:]
				//if writeOffset+len(restBuf) > lenBuf {
				//	c.diag.Error("full", errors.New("full"))
				//	break
				//}

				copy(handleRecvBuf[writeOffset:writeOffset+len(restBuf)], restBuf)
				writeOffset += len(restBuf)
				break
			} else {
				// 找到结束字符，结合缓冲进行处理
				targetBuf := append(handleRecvBuf[0:writeOffset], buf[readOffset:readOffset+index]...)
				fmt.Println(targetBuf)
				//err := c.handlePackageOPPayload(targetBuf)
				//if err != nil {
				//	//数据需要丢弃
				//	//c.diag.Error("handlePackageOPPayload Error", err)
				//	//c.diag.Debug(fmt.Sprintf("procHandleRecv Raw Msg:%s", string(buf)))
				//	//c.diag.Debug(fmt.Sprintf("procHandleRecv Rest Msg:%s writeOffset:%d readOffset:%d index:%d", string(handleRecvBuf), writeOffset, readOffset, index))
				//}

				writeOffset = 0
				readOffset += index + 1
			}
		}
	}
}

func Test_ParseResult001(t *testing.T) {
	result001 := "010000025103程序1                  04                         0500010600107040800000090004100000110120131140151161171181191200000008192210000002200010023000000240000032500100260020027001502800000290000030000003100000320003300034000350000003600000037000000380000003900000040000000410000000219420000043000004417F29015      452020-12-10:14:38:28462020-11-20:06:17:2147拧紧程序             481490150                         51                         52                         53    54000000550000000000560157015800000300000"

	resultdata001 := ResultData{}
	err := ascii.Unmarshal(result001, &resultdata001)
	assert.Nil(t, err)
}
