package utils

import "strconv"

func ParseInt(s string, def int) int {
    if r, err := strconv.Atoi(s); err == nil {
        return r
    }

    return def
}

func ParseNullableInt(s string) *int {
    if s == "" {
        return nil
    }
    if r, err := strconv.Atoi(s); err == nil {
        return &r
    }

    return nil
}