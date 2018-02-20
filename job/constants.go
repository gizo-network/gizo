package job

import "errors"

var (
	ErrExecNotFound           = errors.New("Exec Not Found")
	ErrInvalidPriority        = errors.New("Invalid priority number")
	ErrRetriesOutsideLimit    = errors.New("Retries outside limit")
	ErrRetryDelayOutsideLimit = errors.New("Retry Delay outside limit")
	ErrExecutionTimeBehind    = errors.New("Execution time is past")
)

const (
	MaxRetries      = 5
	MaxRetryDelay   = 120 //! 2 minutes
	DefaultRetries  = 0
	DefaultPriority = NORMAL
)

//! priorities
const (
	HIGH   = 3
	MEDIUM = 2
	LOW    = 1
	NORMAL = 0 //! default
)

//TODO: add errors
//! statuses
const (
	RUNNING     = "running"    //job executed
	FINISHED    = "finished"   //job done
	RETRYING    = "retrying"   //job retrying
	DISPATHCHED = "dispatched" //job dispatched to worker
	STARTED     = "started"    //job received by dispatcher (prior to dispatch)
)
