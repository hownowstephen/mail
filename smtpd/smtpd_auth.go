package smtpd

import (
    "fmt"
    "strings"
)

type Auth struct {
    Mechanisms map[string]AuthExtension
}

func NewAuth() *Auth {
    return &Auth{
        Mechanisms: make(map[string]AuthExtension),
    }
}

func (a *Auth) Handle(c *SMTPConn, args string) error {

    mech := strings.SplitN(args, " ", 2)

    if m, ok := a.Mechanisms[mech[0]]; ok {
        return m.Handle(c, mech[1])
    } else {
        return &AuthError{500, fmt.Errorf("AUTH mechanism %v not available", mech[0])}
    }
}

func (a *Auth) EHLO() string {
    var mechanisms []string
    for m := range a.Mechanisms {
        mechanisms = append(mechanisms, m)
    }
    return strings.Join(mechanisms, " ")
}

func (a *Auth) Extend(mechanism string, extension AuthExtension) error {
    if _, ok := a.Mechanisms[mechanism]; ok {
        return fmt.Errorf("AUTH mechanism %v is already implemented", mechanism)
    }
    a.Mechanisms[mechanism] = extension
    return nil
}

// http://tools.ietf.org/html/rfc4422#section-3.1
// https://en.wikipedia.org/wiki/Simple_Authentication_and_Security_Layer
type AuthExtension interface {
    Handle(*SMTPConn, string) error
}

type AuthPlain struct{}

func (a *AuthPlain) Handle(conn *SMTPConn, params string) error {

    if strings.TrimSpace(params) == "" {
        conn.WriteSMTP(334, "")
        conn.ReadSMTP()
        return nil
    }

    return nil
}

type AuthError struct {
    code int
    err  error
}

func (a *AuthError) Code() int {
    return a.code
}

func (a *AuthError) Error() string {
    return a.err.Error()
}
