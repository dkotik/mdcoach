package mdcoach

import (
	"bytes"

	"github.com/spf13/viper"
)

func setFirstStringValue(key string, m map[string]interface{}, alts ...string) {
	var ok bool
	var s, alt string
	for _, alt = range alts {
		if s, ok = m[alt].(string); ok {
			m[key] = s
			return
		}
	}
}

// ParseFrontMatter parse and normalize yaml data into a map.
func ParseFrontMatter(in []byte) map[string]interface{} {
	result := make(map[string]interface{})
	context := viper.New()
	context.SetConfigType(`yaml`)
	context.ReadConfig(bytes.NewReader(in))
	context.Unmarshal(&result)

	var ok bool
	// TODO: check if case even matters for those or not.
	// TODO: document this feature!
	if _, ok = result[`title`]; !ok {
		setFirstStringValue(`title`, result, `Заглавие`, `заглавие`, `Заголовок`, `заголовок`)
	}
	if _, ok = result[`description`]; !ok {
		setFirstStringValue(`description`, result, `Описание`, `описание`)
	}
	if _, ok = result[`keywords`]; !ok {
		setFirstStringValue(`keywords`, result, `Тэги`, `тэги`)
	}
	if _, ok = result[`author`]; !ok {
		setFirstStringValue(`author`, result, `by`, `Автор`, `aвтор`)
	}
	if _, ok = result[`edited`]; !ok {
		setFirstStringValue(`edited`, result, `changed`, `Редактировано`, `редактировано`)
	}
	if _, ok = result[`created`]; !ok {
		setFirstStringValue(`created`, result, `date`, `Создано`, `создано`, `Дата`, `дата`)
	}
	return result
}
