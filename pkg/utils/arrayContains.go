package utils

import "strings"

func ArrayContains(list []string, item string) bool {
    for _, v := range list {
        if strings.Contains(v, item) {
            return true
        }
    }
    return false
}
