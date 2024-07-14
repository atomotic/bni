package main

import "encoding/xml"

type Collection struct {
	XMLName        xml.Name `xml:"collection"`
	Text           string   `xml:",chardata"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Rec            []Rec    `xml:"rec"`
}

type Rec struct {
	Text string `xml:",chardata"`
	Lab  string `xml:"lab"`
	Cf   []struct {
		Text string `xml:",chardata"`
		T    string `xml:"t,attr"`
	} `xml:"cf"`
	Df []struct {
		Text string `xml:",chardata"`
		T    string `xml:"t,attr"`
		I1   string `xml:"i1,attr"`
		I2   string `xml:"i2,attr"`
		Sf   []struct {
			Text string `xml:",chardata"`
			C    string `xml:"c,attr"`
		} `xml:"sf"`
		S1 []struct {
			Text string `xml:",chardata"`
			Cf   struct {
				Text string `xml:",chardata"`
				T    string `xml:"t,attr"`
			} `xml:"cf"`
			Df struct {
				Text string `xml:",chardata"`
				T    string `xml:"t,attr"`
				I1   string `xml:"i1,attr"`
				I2   string `xml:"i2,attr"`
				Sf   []struct {
					Text string `xml:",chardata"`
					C    string `xml:"c,attr"`
				} `xml:"sf"`
			} `xml:"df"`
		} `xml:"s1"`
	} `xml:"df"`
}

func (r *Rec) ID() string {
	for _, cf := range r.Cf {
		if cf.T == "001" {
			return cf.Text
		}
	}
	return ""
}

func (r *Rec) ISBN() string {
	for _, df := range r.Df {
		if df.T == "010" {
			for _, sf := range df.Sf {
				if sf.C == "a" {
					return sf.Text
				}
			}
		}
	}
	return ""
}

func (r *Rec) Title() string {
	for _, df := range r.Df {
		if df.T == "200" {
			for _, sf := range df.Sf {
				if sf.C == "a" {
					return sf.Text
				}
			}
		}
	}
	return ""
}
