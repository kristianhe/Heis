package common

const (
	DEFAULT_IP   = "127.0.0.1"
	ON           = 1
	OFF          = 0
	TRUE         = 1
	FALSE        = 0
	DISCONNECTED = 0
	CONNECTED    = 1
	EQUALS       = 0
	ELEVATORS    = 2
	INVALID      = -1
	STOP         = 0
	UP           = 2
	DOWN         = 1
	FLOORS       = 4
	BUTTONS      = 3

	MOTOR_UP 	 = 1
	MOTOR_DOWN 	 =-1
	MOTOR_STOP   = 0

	FLOOR_FIRST  = 0
	FLOOR_SECOND = 1
	FLOOR_THIRD  = 2
	FLOOR_LAST   = 3

	BUTTON_INVALID     = -1
	BUTTON_UP          = 0
	BUTTON_DOWN        = 1
	BUTTON_INSIDE      = 2
	BUTTON_OUTSIDE     = 3
	BUTTON_STOP        = 4
	BUTTON_OBSTRUCTION = 5

	STATE_INVALID     = -1
	STATE_STARTUP     = 0
	STATE_IDLE        = 1
	STATE_RUNNING     = 2
	STATE_EMERGENCY   = 3
	STATE_DOOR_OPEN   = 4
	STATE_DOOR_CLOSED = 5

	MESSAGE_INVALID      = -1
	MESSAGE_HEARTBEAT    = 0
	MESSAGE_STATUS       = 1
	MESSAGE_ORDER        = 2
	MESSAGE_FULFILLED    = 3
	MESSAGE_ORDERS       = 4
	MESSAGE_REQUEST      = 5
	MESSAGE_REPRIORITIZE = 6
)

type MotorDirection int

const (
	DIR_UP MotorDirection 	= 1
	DIR_DOWN MotorDirection =-1
	DIR_STOP MotorDirection			= 0
)