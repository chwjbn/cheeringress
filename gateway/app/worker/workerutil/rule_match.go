package workerutil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func GetMatchContext(ctx *gin.Context) map[string]interface{} {

	xDataMap := make(map[string]interface{})

	xDataMap["s_header_url"] = ctx.Request.RequestURI
	xDataMap["s_header_referer"] = ctx.Request.Referer()
	xDataMap["s_header_useragent"] = ctx.Request.UserAgent()
	xDataMap["s_header_cookie"] = ctx.Request.Header.Get("Cookie")
	xDataMap["s_header_x_tenant"] = ctx.Request.Header.Get("X-Tenant")
	xDataMap["s_client_ip"] = ctx.ClientIP()

	return xDataMap

}

//字符串规则匹配
func StringRuleMatched(srcData string, ruleMatchOp string, ruleMatchVal string) bool {

	xRet := false

	if strings.EqualFold(ruleMatchOp, "eq") {
		if strings.EqualFold(srcData, ruleMatchVal) {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "contain") {
		if strings.Contains(srcData, ruleMatchVal) {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "regex") {
		xIsRegexMatch, _ := regexp.MatchString(ruleMatchVal, srcData)
		if xIsRegexMatch {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "in") {
		if strings.Contains(fmt.Sprintf("|%s|", ruleMatchVal), fmt.Sprintf("|%s|", srcData)) {
			xRet = true
			return xRet
		}
	}

	//反向逻辑
	if strings.EqualFold(ruleMatchOp, "neq") {
		if !strings.EqualFold(srcData, ruleMatchVal) {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "notcontain") {
		if !strings.Contains(srcData, ruleMatchVal) {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "notregex") {
		xIsRegexMatch, _ := regexp.MatchString(ruleMatchVal, srcData)
		if !xIsRegexMatch {
			xRet = true
			return xRet
		}
	}

	if strings.EqualFold(ruleMatchOp, "notin") {
		if !strings.Contains(fmt.Sprintf("|%s|", ruleMatchVal), fmt.Sprintf("|%s|", srcData)) {
			xRet = true
			return xRet
		}
	}

	return xRet

}
