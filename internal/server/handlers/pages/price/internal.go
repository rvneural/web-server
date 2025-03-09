package price

import (
	"strconv"
	"time"
)

type ParcedPrice struct {
	Date  time.Time `json:"date"`
	Price uint64    `json:"price"`
}

type ParcedVariation struct {
	Name   string         `json:"name"`
	URL    string         `json:"url"`
	Prices []*ParcedPrice `json:"price"`
}

type ParsedProductElement struct {
	Name      string             `json:"name"`
	Variation []*ParcedVariation `json:"variants"`
}

type ParsedProduct struct {
	Name     string                  `json:"name"`
	Elements []*ParsedProductElement `json:"elements"`
}

func (p *Page) jsonToHTMLTable(data *ParsedProduct) string {
	if len(data.Elements) == 0 {
		return ""
	}
	if len(data.Elements[0].Variation[0].Prices) == 0 {
		return ""
	}

	var table string = "<table>\n"
	table += "<caption>" + "Данные собраны без учета скидок площадки" + "</caption>\n"
	table += "<thead>\n"
	table += "<tr>"
	table += "<td colspan=\"2\" class=\"caption\">" + data.Name + "</td>"
	for _, v := range data.Elements[0].Variation[0].Prices {
		table += "<th class=\"table-date\">" + v.Date.Format("02.01") + "</th>\n"
	}
	table += "</tr>\n"
	table += "</thead>\n"
	col := 1
	for _, v := range data.Elements {
		if col%2 == 0 {
			table += "<tr class=\"tbg\">\n"
		} else {
			table += "<tr>\n"
		}
		col += 1
		rowspan := strconv.Itoa(len(v.Variation))
		table += "<th rowspan=\" " + rowspan + "\" class=\"group\">" + v.Name + "</th>\n"
		for i, s := range v.Variation {
			if i != 0 {
				if col%2 == 0 {
					table += "<tr class=\"tbg\">\n"
				} else {
					table += "<tr>\n"
				}
				col += 1
			}
			table += "<td class=\"variation\">" + s.Name + "</td>\n"
			for _, p := range s.Prices {
				table += "<td>" + strconv.FormatUint(p.Price, 10) + "</td>\n"
			}
			table += "</tr>\n"
		}
	}
	table += "</table>\n"
	return table
}

type JSON struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Models struct {
	Models []JSON `json:"models"`
}
