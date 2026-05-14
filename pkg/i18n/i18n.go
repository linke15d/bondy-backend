package i18n

var messages = map[string]map[string]string{
	"email_invalid": {
		"zh": "邮箱格式不正确",
		"en": "Invalid email format",
		"ja": "メールアドレスの形式が正しくありません",
	},
	"password_too_short": {
		"zh": "密码至少6位",
		"en": "Password must be at least 6 characters",
		"ja": "パスワードは6文字以上必要です",
	},
	"nickname_required": {
		"zh": "昵称不能为空",
		"en": "Nickname is required",
		"ja": "ニックネームは必須です",
	},
	"gender_invalid": {
		"zh": "性别只能选择男、女或其他",
		"en": "Gender must be male, female or other",
		"ja": "性別はmale、female、otherのいずれかを選択してください",
	},
	"email_registered": {
		"zh": "该邮箱已被注册",
		"en": "Email already registered",
		"ja": "このメールアドレスはすでに登録されています",
	},
	"register_failed": {
		"zh": "注册失败，请重试",
		"en": "Registration failed, please try again",
		"ja": "登録に失敗しました。もう一度お試しください",
	},
	"server_error": {
		"zh": "服务器错误",
		"en": "Server error",
		"ja": "サーバーエラー",
	},
	"invalid_params": {
		"zh": "参数格式错误",
		"en": "Invalid parameters",
		"ja": "パラメータの形式が正しくありません",
	},
}

// Get 根据 key 和语言返回对应提示，找不到则返回英文，再找不到返回 key 本身
func Get(key, lang string) string {
	if msgs, ok := messages[key]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
		// fallback 英文
		if msg, ok := msgs["en"]; ok {
			return msg
		}
	}
	return key
}

// ParseLang 从 Accept-Language header 提取主语言
func ParseLang(acceptLanguage string) string {
	if acceptLanguage == "" {
		return "zh"
	}
	// 取第一个语言标签，如 "zh,zh;q=0.9,en;q=0.8" → "zh"
	lang := acceptLanguage
	if idx := len(acceptLanguage); idx > 0 {
		for i, c := range acceptLanguage {
			if c == ',' || c == ';' {
				lang = acceptLanguage[:i]
				break
			}
		}
	}
	supported := map[string]bool{"zh": true, "en": true, "ja": true}
	if supported[lang] {
		return lang
	}
	return "en"
}
