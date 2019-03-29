# Elevated students
Elevator project in TTK4145 Real-time programming.
## Module overview
Cases:         Polling, heartbeat, safemode and more      
Common:        Channels, constant values and message formats    
Control:       Communication with hardware              
Network:       Master-slave communication and UDP broadcasting  
Orders:        Order handling   
Spawn:         Start- and backup procedures   
State machine: State of the elevators   

## Known deficiencies
Elevators has problems with redistributing an order when beeing terminated after having taken the order. Probably has something to do with the spawn routine. Did not work on FAT.

Is not compatible with the elevator simulator due to the master-slave setup. Works splendid with real hardware though.
