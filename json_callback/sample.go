package main

import (
	"encoding/json"
	"fmt"
	"wbapi_golang/wbjsonmsgcrypt"
)

type MsgContent struct {
	ToUsername   string `json:"ToUserName"`
	FromUsername string `json:"FromUserName"`
	CreateTime   uint32 `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
	Msgid        uint64 `json:"MsgId"`
	Agentid      uint32 `json:"AgentId"`
}

func main() {
	token := "hJqcu3uJ9Tn2gXPmxx2w9kkCkCE2EPYo"
	receiverId := "ww1436e0e65a779aee"
	encodingAeskey := "6qkdMrq68nTKduznJYO1A37W2oEgpkMUvkttRToqhUt"
	wbcpt := wbjsonmsgcrypt.NewWBMsgCrypt(token, encodingAeskey, receiverId, wbjsonmsgcrypt.JsonType)
	/*
	   	------------使用示例一：验证回调URL---------------
	   	*企业开启回调模式时，微伴助手会向验证url发送一个get请求
	   	假设点击验证时，企业收到类似请求：
	   	* GET receive_weiban_message?msg_signature=8ddba101fe5f9404e6b28e2d6cfb7565c7572d09&timestamp=1639731768&nonce=64390332&echostr=505deLub74lMxyhj6SZk55rjDRxM7lEzVa/8ZFdOUWkkv9kFpNI%2B9c3YLRLDjL%2B1GEIWpvuQSFw2ck2zb0QcAQ%3D%3D
	   	* HTTP/1.1 Host: weibanzhushou.com

	   	接收到该请求时，企业应
	        1.解析出Get请求的参数，包括消息体签名(msg_signature)，时间戳(timestamp)，随机数字串(nonce)以及微伴助手推送过来的随机加密字符串(echostr),
	        这一步注意作URL解码。
	        2.验证消息体签名的正确性
	        3. 解密出echostr原文，将原文当作Get请求的response，返回给微伴助手
	        第2，3步可以用微伴助手提供的库函数VerifyURL来实现。

	*/
	// 解析出url上的参数值如下：
	// verifyMsgSign := HttpUtils.ParseUrl("msg_signature")
	verifyMsgSign := "8ddba101fe5f9404e6b28e2d6cfb7565c7572d09"
	// verifyTimestamp := HttpUtils.ParseUrl("timestamp")
	verifyTimestamp := "1639731768"
	// verifyNonce := HttpUtils.ParseUrl("nonce")
	verifyNonce := "64390332"
	// verifyEchoStr := HttpUtils.ParseUrl("echoStr")
	verifyEchoStr := "505deLub74lMxyhj6SZk55rjDRxM7lEzVa/8ZFdOUWkkv9kFpNI+9c3YLRLDjL+1GEIWpvuQSFw2ck2zb0QcAQ=="
	echoStr, cryptErr := wbcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if nil != cryptErr {
		fmt.Println("verifyUrl fail", cryptErr)
	}
	fmt.Println("verifyUrl success echoStr", string(echoStr))
	// 验证URL成功，将sEchoStr返回
	// HttpUtils.SetResponse(sEchoStr)

	/*
		------------使用示例二：对用户回复的消息解密---------------
		用户回复消息或者点击事件响应时，企业会收到回调消息，此消息是经过微伴助手加密之后的密文以post形式发送给企业，密文格式请参考官方文档
		假设企业收到微伴助手的回调消息如下：
		POST /receive_weiban_message?msg_signature=cee361999ae632f3e4b9f153475930ce0903b7b7&timestamp=1639732244&nonce=741425964 HTTP/1.1
		Host: weibanzhushou.com
		Content-Length: 613
		{
			"corp_id": "111111111",
			"app_id": "app-114514",
			"encrypt": "3bVgLZLP6TtC8U5qbXxHoq6bL2ZZyarp0s5lyp5RwqNY6E3mhmCU6sb2UlizgL8sHwFqqEtgpOC1OuHBKXF987LB/l8HSRaRHfYseZ5/9QWOJftb/nnI2t6lsEB5aIDkOv3RuvdOYEsvoPy45v6eCbe8sR3WMvy5YDjBiauLWw9us9xvGFx/aRV3DTjPzf8J6H5u8D5339MUP0I+nM11Mhe9pIFnnWMY6Cbd6g9fpwXKalRWBU4N6LXF2/eMypBuVsvJRAncvergZJoBj3svtc2LbaI70qod/G5OmEtbco5A6BXcoaN9HlszmvYH3XIUS8PyRqBOTaELLgIlSIAtzJEX8J7ZjK65LfgVYY2AZzvxxXkQd/mWHhh90M1HPgGQrCYQSkuRb1jGHgKf4HtWgqFzTY50+YwOmrNY4DAun6jrbLwKweMWNMU7g5Qd3YuEklwnteftmL9iAJxJRnx+EihOR+DZjqWibfRnuhPQLuHkMXLvQIzkD+LHDdOYOaWJ3kXvUwt9BGgLKP9QD/kow6nc6xbnUUI+ELxWo65P9PzrJ0otkh1r59rZiwaUC4FFl8W75dh3EGp8dg3Y2adtfekz3Ok9OxMW07HcVoedYdaapUQh8gozu5LC3VZ3Tmmv8TMY/XhWzj8mWIpI6eE57pWTZT3MANhILBmm/n0HmocqIXALjcllklMWWm0X6Z+HEOaf4DiSGk8bRgl9jgifbG0B+s5hV5+1qBMWLaGBLi4MWClldW6wk1YdonlD3fKl0Uxma0eQSVGFlTMezdcJi+92higQevlEC+xanZ3E4m0="
		}

		企业收到post请求之后应该：
			1.解析出url上的参数，包括消息体签名(msg_signature)，时间戳(timestamp)以及随机数字串(nonce)
			2.验证消息体签名的正确性。
			3.将post请求的数据进行json解析，并将"Encrypt"标签的内容进行解密，解密出来的明文即是用户回复消息的明文，明文格式请参考官方文档
			第2，3步可以用微伴助手提供的库函数DecryptMsg来实现。
	*/

	// reqMsgSign := HttpUtils.ParseUrl("msg_signature")
	reqMsgSign := "cee361999ae632f3e4b9f153475930ce0903b7b7"
	// reqTimestamp := HttpUtils.ParseUrl("timestamp")
	reqTimestamp := "1639732244"
	// reqNonce := HttpUtils.ParseUrl("nonce")
	reqNonce := "741425964"
	// post请求的密文数据
	// reqData = HttpUtils.PostData()

	reqData := []byte(`{"corp_id": "111111111", "app_id": "app-114514", "encrypt": "3bVgLZLP6TtC8U5qbXxHoq6bL2ZZyarp0s5lyp5RwqNY6E3mhmCU6sb2UlizgL8sHwFqqEtgpOC1OuHBKXF987LB/l8HSRaRHfYseZ5/9QWOJftb/nnI2t6lsEB5aIDkOv3RuvdOYEsvoPy45v6eCbe8sR3WMvy5YDjBiauLWw9us9xvGFx/aRV3DTjPzf8J6H5u8D5339MUP0I+nM11Mhe9pIFnnWMY6Cbd6g9fpwXKalRWBU4N6LXF2/eMypBuVsvJRAncvergZJoBj3svtc2LbaI70qod/G5OmEtbco5A6BXcoaN9HlszmvYH3XIUS8PyRqBOTaELLgIlSIAtzJEX8J7ZjK65LfgVYY2AZzvxxXkQd/mWHhh90M1HPgGQrCYQSkuRb1jGHgKf4HtWgqFzTY50+YwOmrNY4DAun6jrbLwKweMWNMU7g5Qd3YuEklwnteftmL9iAJxJRnx+EihOR+DZjqWibfRnuhPQLuHkMXLvQIzkD+LHDdOYOaWJ3kXvUwt9BGgLKP9QD/kow6nc6xbnUUI+ELxWo65P9PzrJ0otkh1r59rZiwaUC4FFl8W75dh3EGp8dg3Y2adtfekz3Ok9OxMW07HcVoedYdaapUQh8gozu5LC3VZ3Tmmv8TMY/XhWzj8mWIpI6eE57pWTZT3MANhILBmm/n0HmocqIXALjcllklMWWm0X6Z+HEOaf4DiSGk8bRgl9jgifbG0B+s5hV5+1qBMWLaGBLi4MWClldW6wk1YdonlD3fKl0Uxma0eQSVGFlTMezdcJi+92higQevlEC+xanZ3E4m0="}`)

	msg, cryptErr := wbcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, reqData)
	if nil != cryptErr {
		fmt.Println("DecryptMsg fail", cryptErr)
	}
	fmt.Println("after decrypt msg: ", string(msg))
	// TODO: 解析出明文json标签的内容进行处理
	// For example:

	var msgContent MsgContent
	err := json.Unmarshal(msg, &msgContent)
	if nil != err {
		fmt.Println("Unmarshal fail", err)
	} else {
		fmt.Println("struct", msgContent)
	}

	/*
		------------使用示例三：企业回复用户消息的加密---------------
		企业被动回复用户的消息也需要进行加密，并且拼接成密文格式的json串。
		假设企业需要回复用户的明文如下：

		{
			"create_time": 1348831860,
			"corp_id": "11111",
			"Content": "this is a test",
			"app_id": "128"
		}

		为了将此段明文回复给用户，企业应：
			1.自己生成时间时间戳(timestamp),随机数字串(nonce)以便生成消息体签名，也可以直接用从微伴助手的post url上解析出的对应值。
			2.将明文加密得到密文。
			3.用密文，步骤1生成的timestamp,nonce和企业在微伴助手设定的token生成消息体签名。
			4.将密文，消息体签名，时间戳，随机数字串拼接成json格式的字符串，发送给企业。
			以上2，3，4步可以用微伴助手提供的库函数EncryptMsg来实现。
	*/
	respData := `{"create_time": 1348831860,"corp_id": "11111","Content": "this is a test","app_id": "128"}`
	encryptMsg, cryptErr := wbcpt.EncryptMsg(respData, reqTimestamp, reqNonce)
	if nil != cryptErr {
		fmt.Println("DecryptMsg fail", cryptErr)
	}

	sEncryptMsg := string(encryptMsg)

	fmt.Println("after encrypt sEncryptMsg: ", sEncryptMsg)
	// 加密成功
	// TODO:
	// HttpUtils.SetResponse(sEncryptMsg)
}
