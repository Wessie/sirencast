package taxtic

import (
	"errors"
	"strings"

	"code.google.com/p/go.text/encoding"
	"code.google.com/p/go.text/encoding/charmap"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/encoding/korean"
	"code.google.com/p/go.text/encoding/simplifiedchinese"
	"code.google.com/p/go.text/encoding/traditionalchinese"
	"code.google.com/p/go.text/encoding/unicode"
	"code.google.com/p/go.text/transform"
)

var (
	ErrUnknownCharset = errors.New("taxtic: unknown charset")
	ErrUnrecoverable  = errors.New("taxtic: unrecoverable error")
)

func Convert(charset, s string) (string, error) {
	charset = strings.Replace(charset, "-", "", -1)
	charset = strings.ToLower(charset)

	if charset == "utf8" {
		return s, nil
	}

	e := Encoding(charset)
	if e == nil {
		return "", ErrUnknownCharset
	}

	d := e.NewDecoder()

	res, _, err := transform.Bytes(d, []byte(s))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func Encoding(charset string) encoding.Encoding {
	charset = strings.Replace(charset, "-", "", -1)
	charset = strings.ToLower(charset)

	switch charset {
	case "sjis", "cp932":
		return japanese.ShiftJIS
	case "eucjp":
		return japanese.EUCJP
	case "iso2022jp":
		return japanese.ISO2022JP

	case "euckr", "cp949":
		return korean.EUCKR
	case "koi8r":
		return charmap.KOI8R
	case "koi8u":
		return charmap.KOI8U

	case "macintosh":
		return charmap.Macintosh
	case "macintoshcyrillic":
		return charmap.MacintoshCyrillic

	case "gb18030":
		return simplifiedchinese.GB18030
	case "gbk", "cp936":
		return simplifiedchinese.GBK
	case "hzgb2312":
		return simplifiedchinese.HZGB2312
	case "big5", "cp950":
		return traditionalchinese.Big5

	case "cp437":
		return charmap.CodePage437
	case "cp866":
		return charmap.CodePage866
	case "iso8859_2":
		return charmap.ISO8859_2
	case "iso8859_3":
		return charmap.ISO8859_3
	case "iso8859_4":
		return charmap.ISO8859_4
	case "iso8859_5":
		return charmap.ISO8859_5
	case "iso8859_6":
		return charmap.ISO8859_6
	case "iso8859_7":
		return charmap.ISO8859_7
	case "iso8859_8":
		return charmap.ISO8859_8
	case "iso8859_10":
		return charmap.ISO8859_10
	case "iso8859_13":
		return charmap.ISO8859_13
	case "iso8859_14":
		return charmap.ISO8859_14
	case "iso8859_15":
		return charmap.ISO8859_15
	case "iso8859_16":
		return charmap.ISO8859_16

	case "windows1250":
		return charmap.Windows1250
	case "windows1251":
		return charmap.Windows1251
	case "windows1252", "latin1":
		return charmap.Windows1252
	case "windows1253":
		return charmap.Windows1253
	case "windows1254":
		return charmap.Windows1254
	case "windows1255":
		return charmap.Windows1255
	case "windows1256":
		return charmap.Windows1256
	case "windows1257":
		return charmap.Windows1257
	case "windows1258":
		return charmap.Windows1258
	case "windows874":
		return charmap.Windows874

	case "utf16be":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "utf16le":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	}
	return nil
}
