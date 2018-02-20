package job

//TODO: add errors
const (
	RUNNING     = "running"    //job executed
	FINISHED    = "finished"   //job done
	RETRYING    = "retrying"   //job retrying
	DISPATHCHED = "dispatched" //job dispatched to worker
	STARTED     = "started"    //job received by dispatcher (prior to dispatch)
)
