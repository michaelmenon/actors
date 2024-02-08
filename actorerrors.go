package actors

type ActorError struct {
	err string
}

// implement error interface
func (ae ActorError) Error() string {
	return ae.err
}

const ACTORHUBGENERROR = "cannot create ActorsHub"
const NILSTOREERROR = "ActorHub cache is not initialized."
const DUPLICATEACTORERROR = "actor with the given id already exist"
const ACTORNOTFOUNDERROR = "actor does not exist"
const ACTORCHANERROR = "actor receive channel does not exist"
const ACTORMAXLIMITERROR = "maximum actors limit reached on a tag"
const HUBALREADYRUNNING = "Hub already running"
const HUBNOTRUNNING = "Hub not runnin"
