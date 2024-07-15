package repoerrs

import "errors"

var ErrUserAlreadyVoted = errors.New("user already voted")
var ErrPollNotExist = errors.New("poll not exist")
