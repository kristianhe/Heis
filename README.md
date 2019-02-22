# Heis
Repository for group 43 to work on the elevator project.

Exercise 3: Uplink Established
==============================

This exercise is part of of a three-stage process:
 - This first exercise is to make you more familiar with using TCP and UDP for communication between processes running on different machines. Do not think too much about code quality here, as the main goal is familiarization with the protocols.
 - Exercise 4 will have you consider the things you have learned about these two protocols, and implement (or find) a network module that you can use in the project. The communication that is necessary for your design should reflect your choice of protocol. This is when you should start thinking more about code quality, because ...
 - Finally, you should be able to use this module as part of your project. Since predicting the future is notoriously difficult, you may find you need to change some functionality. But if the module has well-defined boundaries (a set of functions, communication channels, or others), you *should* be able to swap out the entire thing if necessary!


Note:
 - You are free to choose any language. Using the same language on the network exercises and the project is recommended, but not required. If you are still in the process of deciding, use this exercise as a language case study.
 - Exactly how you do communication for the project is up to you, so if you want to venture out into the land of libraries you should make sure that the library satisfies all your requirements. You should also check the license.

Practical tips:
 - Sharing a socket between threads should not be a problem, although reading from a socket in two threads will probably mean that only one of the threads get the message. If you are using blocking sockets, you could create a "receiving"-thread for each socket. Alternatively, you can use socket sets and the [`select()`](http://en.wikipedia.org/wiki/Select_%28Unix%29) function (or its equivalent).
 - Be nice to the network: Put some amount of `sleep()` or equivalent in the loops that send messages. The network at the lab will be shut off if IT finds a DDOS-esque level of traffic. Yes, this has happened before. Several times.
 - You can find [some pseudocode here](resources.md) to get you started.


Part 1: UDP
-----------

We have set up a server on the real time lab that you're going to communicate with in this exercise. If you cannot connect to it, it might be down. Ask a student assistant to turn it on for you.

### Receiving UDP packets, and finding the server IP:
The server broadcasts it's own ip on port `30000`, listen for messages on this port to find the IP. You should write down the IP address as you will need it for again later in the exercise.

### Sending UDP packets:
The server will respond to any message you send to it. Try sending a message to the server ip on port `20000 + n` where `n` is the number of the workspace you're sitting at. Listen to messages from the server and print them to a terminal to confirm that the messages are in fact beeing responded to.

- The server will act the same way if you send a broadcast (`#.#.#.255` or `255.255.255.255`) instead of sending directly to the server.
  - If you use broadcast, the messages will loop back to you! The server prefixes its reply with "You said: ", so you can tell if you are getting a reply from the server or if you are just listening to yourself.
- You are free to mess with the people around you by using the same port as them, but they may not appreciate it.
- It may be easier to use two sockets: one for sending and one for receiving. You might also find it is easier if these two are separated into their own threads.


Part 2: TCP
-----------

There are three common ways to format a message: (Though there are probably others)
 - 1: Always send fixed-sized messages
 - 2: Send the message size with each message (as part of a fixed-size header)
 - 3: Use some kind of marker to signify the end of a message

The TCP server can send you two of these:
 - Fixed size messages of size `1024`, if you connect to port `34933`
 - Delimited messages that use `\0` as the marker, if you connect to port `33546`

The server will read until it encounters the first `\0`, regardless. Strings in most programming languages are already null-terminated, but in case they aren't you will have to append your own null character.



### Connecting:
The IP address of the TCP server will be the same as the address the UDP server as spammed out on port 30000. Connect to the TCP server, using port `34933` for fixed-size messages, or port `33546` for 0-terminated messages. It will reply back with a welcome-message when you connect. 

The server will then echo anything you say to it back to you (as long as your message ends with `'\0'`). Try sending and receiving a few messages.

### Accepting connections:
Tell the server to connect back to you, by sending a message of the form `Connect to: #.#.#.#:#\0` (IP of your machine and port you are listening to). You can find your own address by running `ifconfig` in the terminal, the first three bytes should be the same as the server's IP.
 
This new connection will behave the same way on the server-side, so you can send messages and receive echoes in the same way as before. You can even have it "Connect to" recursively (but please be nice... And no, the server will refuse requests to connect to itself).
 



Exercise 4: From Prototype to Production
========================================

1. Don't overengineer.
2. Always design properly.
3. Minor detail change will ruin your perfect model.
4. Always prototype first.
5. You will only know what you didn't know on the second iteration.
6. There is only budget/time for the first version.
7. The prototype will become production.


## Network Module

This exercise aims to produce the "v0.1" of the network module used in your project and at the same time prepare you for the network part on the design review.


You should start by taking a look back at [the beginning of Exercise 1](https://github.com/TTK4145/Exercise1/blob/master/Part1/README.md), and reevaluate them in the light of what you have learned about network programming (and - if applicable - concurrency control). At the same time you might want to look into what kind of libraries that already exist for your chosen language.

 - [C network module](https://github.com/TTK4145/Network-c)
 - [D network module](https://github.com/TTK4145/Network-D)
 - [Go network module](https://github.com/TTK4145/Network-go)
 - [Rust network module](https://github.com/edvardsp/network-rust)
 - [Distributed Erlang](http://erlang.org/doc/reference_manual/distributed.html)
 
By the end of this exercise, you should be able to send some data structure (struct, record, etc) from one machine to another. How you achieve this (in terms of network topology, protocol, serialization) does not matter. The key here is *abstraction*.  

Don't forget that this module should *simplify* the interface between machines: Creating and handling sockets in all routines that need to communicate with the outside world is possible, but is likely to be unwieldy and unmaintainable. We want to encapsulate all the necessary functionality in a single module, so we have a single decoupled component where we can say "This module sends our data over the network". This will almost always be preferable, but above all else: *Think about what best suits your particular design*.


### Design questions

To get you started with designing a network module and/or the application that uses it, try to find answers to the questions below:

 ##### Guarantees about elevators:
What should happen if one of the nodes loses network?
> All nodes check in periodically. If one node doesn't check in, it means the node has lost network or lost power. If it lost network, it will continue to function locally (within the elevator panel), but not respond to instructions from the floor panel. The other elevators will continue to work, and handle all tasks from the floor panel.

What should happen if one of the computers loses power for a brief moment?
> The data should not be compromised by the power failure in one computer. Data should be exchanged between all nodes often (message passing), so that all nodes carries an updated version of the shared data. The node with a power loss should be able to download the new data from its friends and carry on working once the power is back on.

What should happen if some unforeseen event causes the elevator to never reach its destination but communication remains intact?
> If the elevator is empty and responding; The others should take over the order after some idle time.
> If the elevator has clients inside; It should prioritize getting to the nearest floor and open the doors.

##### Guarantees about orders:
Do all your nodes need to "agree" on an order for it to be accepted? In that case, how is a faulty node handled?
> Say that all nodes are separated; an elevator should respond to an order even if the other nodes doesn't have the same order, or communication is down.

How can you be sure that at least as many nodes as needs to agree on the order in fact agrees on the order?
> Given that more than one elevator needs to agree; there should be a check on each order where all nodes can "agree" on it. Only if The required number of nodes agree, the order will be added to the list of executable orders.

Do you share the entire state of the current orders, or just the changes as they occur?
> The entire state should be shared. That way any nodes off line can receive the full state of the system once they re-engage.

For either one: What should happen when an elevator re-joins after having been off line?
> See above.

##### Topology:
What kind of network topology do you want to implement? Peer to peer? Master slave? Circle?
> We want a broadcast-like topology which enables the possibility for each elevator to broadcast information at a certain rate. In this way we can label the messages with ID's, much like a CAN-bus network. Practically, this means that all the elvators will listen at all times, and choose to receive the information based on the message ID.

In the case of a master-slave configuration: Do you have only one program, or two (a "master" executable and a "slave")? Is a slave becoming a master a part of the network module?
> If we avoid a master-slave configuration we also avoid the problem with power loss on the master computer which cripples the whole system.
    
##### Technical implementation:
If you are using TCP: How do you know who connects to who?
> Are going to use UDP with some TCP functionality (some kind of acknowledge and retransmission).

Do you need an initialization phase to set up all the connections?
> Probably yes.

Will you be using blocking sockets & many threads, or nonblocking sockets and select()?
> We will be using for-select-loops with some threads.

Do you want to build the necessary reliability into the module, or handle that at a higher level?
> We build the modules as simples as possible, and handle realiability and errors at a higher level.

How will you pack and unpack (serialize) data?
Do you use structs, classes, tuples, lists, ...?
> The simplest data: Lists.
JSON, XML, or just plain strings?
> We will use plain strings. The go language has many commands that can manipulate strings in a number of ways. This might
be useful for us when working with the data. 
Is serialization a part of the network module?
> Surely it should be. The network module should handle all that has to do with data transmission.

Is detection (and handling) of things like lost messages or lost nodes a part of the network module?
> Yes it should be. Again we leave to the network module to handle all things regarding message transmission.

##### Protocols:
TCP gives you a data stream that is guaranteed to arrive in the same order as it was sent in (or not at all)  
UDP might reorder the packets you send into the network  
TCP will resend packets if they're dropped (at least until the socket times out)  
UDP may drop packets as it pleases  
TCP requires that you to set up a connection, so you will have to know who connects to who  
UDP doesn't need a connection, and even allows broadcasting  
(You're also allowed to use any other network library or language feature you may desire)  
> We are going to use UDP.

 
### Running from another computer

In order to test a network module, you will have to run your code from multiple machines at once. The best way to do this is to log in remotely. Remember to be nice the people sitting at that computer (don't mess with their files, and so on), and try to avoid using the same ports as them.

 - Logging in:
   - `ssh username@#.#.#.#` where #.#.#.# is the remote IP
 - Copying files between machines:
   - `scp source destination`, with optional flag `-r` for recursive copy (folders)
   - Examples:
     - Copying files *to* remote: `scp -r fileOrFolderAtThisMachine username@#.#.#.#:fileOrFolderAtOtherMachine`
     - Copying files *from* remote: `scp -r username@#.#.#.#:fileOrFolderAtOtherMachine fileOrFolderAtThisMachine`
129.241.229.181
    
*If you are scripting something to automate any part of this process, remember to **not** include the login password in any files you upload to github (or anywhere else for that matter)*

## Extracurricular

[The Night Watch](https://web.archive.org/web/20140214100538/http://research.microsoft.com/en-us/people/mickens/thenightwatch.pdf)
*"Systems people discover bugs by waking up and discovering that their first-born children are missing and "ETIMEDOUT" has been written in blood on the wall."*

[The case of the 500-mile email](http://www.ibiblio.org/harris/500milemail.html)
*"We can't send mail farther than 500 miles from here," he repeated. "A little bit more, actually. Call it 520 miles. But no farther."*

[21 Nested Callbacks](http://blog.michellebu.com/2013/03/21-nested-callbacks/)
*"I gathered from these exchanges that programmers have a perpetual competition to see who can claim the most things as 'simple.'"*


Exercise 5: Call for transport
==============================

The elevator hardware on the lab is controlled via a National Instruments PCI Digital I/O device, using the Comedi driver. An "elevator" abstraction has been created, that exposes a few functions in C that lets us use this I/O card. However, using this presents a few challenges for a project like this:

 - Calling C code from other programming languages is sometimes a bit of a hassle
 - The driver only works on Linux, which you might not use when working elsewhere than the lab
 - Very few of you have an elevator
 
In order to alleviate the lack of elevators, a simulator was created. In order to use the simulator, you need to see what it does and input "button presses" to it, which means it has to run in a separate terminal window. Then in order to a) make interfacing with the simulator and the real elevator as similar as possible, and b) eliminate the need to call C code, both the simulator and the hardware elevator expose a network-based interface using TCP. In this way, you can seamlessly swap between the simulator and hardware elevators.

This means we have a simple client-server structure to the elevator:

 - Two possible servers:
   - [The Elevator Server](https://github.com/TTK4145/elevator-server)
   - [The simulator](https://github.com/TTK4145/Simulator-v2)
 - [Language-specific clients](https://github.com/TTK4145?q=driver)
   - Choose the one you need for the language you are using on the project
   - (If none exist for your language, ask for help and we'll add it once it works)

You may want to modify the client end of the driver, or possibly create your own from scratch. There is no particular requirement or recommendation involved here.

*(The [low-level C driver](/driver) for the elevator hardware is included in this repository for completeness, but using it is not recommended.)*
    
Up and down
-----------

Download the driver (for the programming language you are doing the project in), and test it on both the hardware elevator and the simulator.

 - Using the hardware elevator
   - Download and run the [elevator server executable](https://github.com/TTK4145/elevator-server/releases/latest). It's likely already installed on the lab computers, try just running `ElevatorServer` from the terminal
   - If an `ElevatorServer` is already running, the new server will not be able to bind to the socket. If you need to kill it, you can do so by calling `pkill ElevatorServer`
 - Using the simulator
   - Download and run the [simulator executable](https://github.com/TTK4145/Simulator-v2/releases/latest)
   - In order to run multiple simulators on the same computer, you will have to change the port on both the simulator (with `--port`) and in the driver (likely in a call to some init-function or in a config file)


Up and away
-----------

The elevator project can be roughly divided into two parts: Distributing the incoming requests (hall and cab calls) to the elevators, and then servicing those requests. At this stage you don't have the functionality implemented for the first part, but the latter part was a project in TTK4235. Since not all of you have taken this course, we'll have to get you up to speed on both the solution to this problem, and the preferred implementation pattern.

The relevant part of that project is documented in [Project-resources/elev_algo](https://github.com/TTK4145/Project-resources/tree/master/elev_algo).

Implementing the "single elevator control" component as a state machine is the preferred pattern, to the point where we might even dare call it the *definitively correct* approach. The details and analysis of this pattern are covered in greater detail in the lectures, but here is the short version on how to follow it:

 - Analysis:
   - Identify the inputs of the system, and consider them as discrete (in time) events
   - Identify the outputs of the system
   - Combine all inputs & outputs and store them (the "state") from one event to the next
     - This creates a combinatorial "explosion" of possible internal state
   - Eliminate combinations that are redundant, impossible (according to reality), and undesirable (according to the specification)
     - This should give you a "minimal representation" of the possible internal state
   - Give names to the new minimal combined-input-output states
     - These typically identify how the system "behaves" when responding to the next event
     - Leave any un-combined data alone
 - Implementation:
   - Create a "handler" for each event (function, message receive-case, etc)
   - Switch/match on the behaviour-state (within each event-handler)
   - Call a (preferably pure) function that computes the next state and output actions
   - Perform the output actions and save the new state
   
You are encouraged to try to trace the analysis steps for the elev_algo code linked above, but I also find that the vastly less rigorous approach of intuition quickly overtakes the methodical one. But for the implementation side you should take a much closer look, especially on why we consider events first then state (as opposed to state-first), and where the divide goes between "code that is directly in the event-handler" and "code that is in a function called by the event handler".

### Doing it yourself

You should now implement some way to control a single elevator, as a part of the elevator project. This is where you get started "for real", so set up your environment (build tools, repository, editor, etc) the way you like it before you begin.

Since you don't have any way to distribute requests yet, you should use the button presses directly. This will have to change later, so keep code quality in mind.

Exercise 6 : Phoenix
====================

Create a program (in any language, on any OS) that uses the process pair technique to print the numbers 1, 2, 3, 4, etc to a terminal window. The program should create its own backup: When the primary is running, only the primary should keep counting, and the backup should do nothing. When the primary dies, the backup should become the new primary, create its own new backup, and keep counting where the dead one left off. Make sure that no numbers are skipped!

You cannot rely on the primary telling the backup when it has died (because it would have to be dead first...). Instead, have the primary broadcast that it is still alive, and have the backup become the primary when a certain number of messages have been missed.

You will need some form of communication between the primary and the backup. Some examples are:

    Network: The simplest is to use UDP on localhost. TCP is also possible, but may be harder (since both endpoints need to be alive).
    IPC, such as POSIX message queues: see msgget() msgsnd() and msgrcv(). With these you can create FIFO message queues.
    Signals: Use signals to interrupt other processes (You are already familiar with some of these, such as SIGSEGV (Segfault) and SIGTERM (Ctrl+C)). There are two custom signals you can use: SIGUSR1 and SIGUSR2. See signal().
        Note for D programmers: SIGUSR is used by the GC.
    Files: The primary writes to a file, and the backup reads it. Either the time-stamp of the file or the contents can be used to detect if the primary is still alive.
    Controlled shared memory: The system functions shmget() and shmat() let processes share memory.

You will also need to spawn the backup somehow. There should be a way to spawn processes or run shell commands in the standard library of your language of choice. The name of the terminal window is OS-dependent:

    Ubuntu: gnome-terminal -x ["commands"]
    Windows: start "title" [program_name]. Note that you must specify a title
    OSX: osascript -e 'tell app "Terminal" to do script ["terminal command"]'

(Linux tip: You can prevent a spawned terminal window from automatically closing by going to Edit -> Profile Preferences -> Title and Command -> When command exits. Windows tip: Use start "title" call [program_name])

Be careful! You don't want to create a chain reaction... If you do, you can use pkill -f program_name (Windows: taskkill /F /IM program_name /T) as a sledgehammer.

In case you want to use this on the project: Usually a program crashes for a reason. Restoring the program to the same state as it died in may cause it to crash in exactly the same way, all over again. How would you prevent this from happning?
Guarantees (Optional)

Make your program print each number once and only once, and demonstrate (a priori, not just through observation of your program) that it will behave this way, regardless of when the primary is killed.
Approval

    Demonstrate the process pair functionality (kill the primary and show that the backup takes over).
    Show the code you created.

Exercise 7 : Backward error recovery

This is the first of a two-part exercise: This time we will look at one way of doing backward error recovery, and next time we will modify this code to use forward error recovery instead. Since we will be needing a rather unique language feature of Ada for our forward error recovery solution, we will be using Ada for both parts, so it is easier to see the similarities and differences between the two approaches.
Practical information.

    Have a look at Intro to Ada
    Use the Ada Reference Manual if you need to look up something.
    Compile your program using gnatmake.
    The code you are completing can be found here.

Desired functionality:

We are modelling a transaction with three participants, where each performs a calculation that is slow when it works correctly, but fails quickly when it fails. The work from each should not be committed (in this case: printed to the standard output) unless all participants succeed. When any participant fails, the work from all the others will need to be reset, and the transaction has to start over from the beginning.
Create the transaction work function

The "work" the participants are doing is adding 10 to a number. Unoriginal, perhaps, but we can use random numbers to have it simulate work that either success or fails. We will call this function Unreliable_Slow_Add.

    A random number generator Gen is defined and seeded for you. Call Random(Gen) to get a random number between 0.0 and 1.0. Compare it with the Error_Rate, and have the function either perform:
        The intended behaviour: Most of the time, the addition takes up to 4-ish seconds. Use delay Duration(d) (where d is a floating-point number) to pause execution for d seconds (You can use Random(Gen) multiplied with a constant as the value for d). Then, add 10 to x and return the value.
        The faulty behaviour: The rest of the time, the operation takes significantly less time (say, up to half a second), but raises an exception instead. A Count_Failed exception has been defined for you. (Note: Ada uses raise, not throw)

Do the transaction work

Now that we have the unreliable slow adder, we need to call it, and also catch the exception it throws.

    The variable we are modifying is called Num, and its previous value is called Prev.
    The structure for exception handling in Ada is begin-exception-end. There is only one exception to catch: Count_Failed. When counting fails, we need to tell the transaction manager that everyone has to revert, by using Manager.Signal_Abort; (already implemented).
    Both in case of success and failure, we need to know what happened to the other participants. The exit protocol lies in the Finished entry of the transaction manager, but is not yet implemented (We'll get back to this in part 3). Call it, and trust your future selves that it will be implemented properly.
    We then ask the manager if we should commit the result. If not, we have to revert to the previous value.

Finish the Manager exit protocol

The exit protocol requires that all participants show up, and all votes are counted. We store the vote using two booleans: Aborted, which is set true by the first participant that aborts the transaction (by calling Signal_Abort), and Should_Commit, which stores this value until the next round starts.

The Transaction_Manager is a protected object. You can find a quick summary of protected objects here. We are also using one additional feature: The Count attribute. We can get the number of tasks blocked on the condition of the Finished entry by using Finished'Count.

    The exit protocol needs to function as a "gate", letting participants through only when all of them have arrived (when Finished'Count = N makes sure no task enters until N of them have arrived).
    When the first participant enters, we open the gate for the remaining participants. When the last one enters, the gate is closed (for next round). Use Finished'Count (or another counting variable) to see how many participants are waiting to enter.
    We also need to update and reset the voting variables. When the first participant enters, it can set Should_Commit. When the last participant enters, it can reset Aborted (this is the reason we store Should_Commit as a separate variable).

Approval.

    Show the student assistants that your program produce the expected results: All participants either add 10, or all participants revert.

Next time

These Questions will be answered next time (so think about them now, before you encounter any spoilers):

    We can have a situation where we have multiple failures in a row, denying any progress from happening. What are the real-time consequences of this?
    When one of the participants fails early, all the other participants are doing futile work, yet we cannot do anything but wait until all participants are finished. How can we solve this?



