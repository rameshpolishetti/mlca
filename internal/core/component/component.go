package component

// Component is container component managed by container agent
type Component interface {
	/*
	* UNKNOWN	initialize()	bootup()
	* UNSATISFIED	resolveDependencies()	buildConfiguration()
	* RESOLVED	activate()	launchComponent()
	* STANDBY	standby()	prepareForActive()
	* ACTIVE	monitor()	watchComponent()
	* RELOAD	reload()
	* RECYCLE	waitingForDependencies()
	* DISABLED	deavtivate()
	 */

	Bootup() bool
	BuildConfiguration() bool
	LaunchComponent() bool
	PrepareForActive() bool
	WatchComponent() bool
}
