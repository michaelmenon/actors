package cmd

type ActorError struct {
	Err string
}

// implement error interface
func (ae ActorError) Error() string {
	return ae.Err
}

const ErrActorHubGen = "cannot create ActorsHub"
const ErrNilStore = "ActorHub cache is not initialized."
const ErrDuplicateActor = "actor with the given id already exist"
const ErrActorNotFound = "actor does not exist"
const ErrActorChannel = "actor receive channel does not exist"
const ErrActorMaxLimit = "maximum actors limit reached on a tag"
const ErrHubRunning = "Hub already running"
const ErrHubNotRunning = "Hub not runnin"
const ErrNetActorCreate = "Cannot create network actor"
const ErrClusterConnection = "Cannot connect to master node"
const NoClusterConnection = "No cluster found"
const ClusterMessageError = "Unable to create cluster message to send"
