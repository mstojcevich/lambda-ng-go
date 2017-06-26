// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package user

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser(in *jlexer.Lexer, out *SessionResult) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "errors":
			if in.IsNull() {
				in.Skip()
				out.Errors = nil
			} else {
				in.Delim('[')
				if out.Errors == nil {
					if !in.IsDelim(']') {
						out.Errors = make([]string, 0, 4)
					} else {
						out.Errors = []string{}
					}
				} else {
					out.Errors = (out.Errors)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Errors = append(out.Errors, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "id":
			out.UserID = int(in.Int())
		case "username":
			out.Username = string(in.String())
		case "api_key":
			out.APIKey = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser(out *jwriter.Writer, in SessionResult) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"errors\":")
	if in.Errors == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in.Errors {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.String(string(v3))
		}
		out.RawByte(']')
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Int(int(in.UserID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"username\":")
	out.String(string(in.Username))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"api_key\":")
	out.String(string(in.APIKey))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SessionResult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SessionResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SessionResult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SessionResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser(l, v)
}
func easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser1(in *jlexer.Lexer, out *RegisterResult) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "errors":
			if in.IsNull() {
				in.Skip()
				out.Errors = nil
			} else {
				in.Delim('[')
				if out.Errors == nil {
					if !in.IsDelim(']') {
						out.Errors = make([]string, 0, 4)
					} else {
						out.Errors = []string{}
					}
				} else {
					out.Errors = (out.Errors)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.Errors = append(out.Errors, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "api_key":
			out.APIKey = string(in.String())
		case "success":
			out.Success = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser1(out *jwriter.Writer, in RegisterResult) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"errors\":")
	if in.Errors == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in.Errors {
			if v5 > 0 {
				out.RawByte(',')
			}
			out.String(string(v6))
		}
		out.RawByte(']')
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"api_key\":")
	out.String(string(in.APIKey))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"success\":")
	out.Bool(bool(in.Success))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegisterResult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegisterResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegisterResult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegisterResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser1(l, v)
}
func easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser2(in *jlexer.Lexer, out *LoginResult) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "errors":
			if in.IsNull() {
				in.Skip()
				out.Errors = nil
			} else {
				in.Delim('[')
				if out.Errors == nil {
					if !in.IsDelim(']') {
						out.Errors = make([]string, 0, 4)
					} else {
						out.Errors = []string{}
					}
				} else {
					out.Errors = (out.Errors)[:0]
				}
				for !in.IsDelim(']') {
					var v7 string
					v7 = string(in.String())
					out.Errors = append(out.Errors, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "api_key":
			out.APIKey = string(in.String())
		case "success":
			out.Success = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser2(out *jwriter.Writer, in LoginResult) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"errors\":")
	if in.Errors == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in.Errors {
			if v8 > 0 {
				out.RawByte(',')
			}
			out.String(string(v9))
		}
		out.RawByte(']')
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"api_key\":")
	out.String(string(in.APIKey))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"success\":")
	out.Bool(bool(in.Success))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LoginResult) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LoginResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson42239ddeEncodeGithubComMstojcevichLambdaNgGoUser2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LoginResult) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LoginResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson42239ddeDecodeGithubComMstojcevichLambdaNgGoUser2(l, v)
}
