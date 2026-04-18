Port Forwarder

Author : Abhijit Rekhi

A lightweight, concurrent TCP relay and static port forwarder written in Go. Designed for Red Team operations to facilitate network pivoting and internal tunneling across restricted network enclaves without requiring external dependencies or interpreters.

## Operational Features 
 * Zero Dependencies: Compiles to a single static binary, ensuring execution compatibility across diverse target environments.
 * High Concurrency: Utilizes Go's lightweight goroutines for non-blocking, asynchronous I/O. Capable of handling thousands of concurrent connections with minimal memory overhead.
 * Leak Protection: Implements size-buffered channels for full-duplex stream synchronization, explicitly preventing goroutine "zombie" leaks during asynchronous TCP cancellation or abrupt client disconnects.
 * Kernel-Optimized Transfers: Uses native io.Copy to leverage underlying OS efficiencies (like splice on Linux) for memory-to-memory byte transfers.
 * Cross-Platform: Native cross-compilation support for Windows, Linux, and macOS target architectures.
   
Compilation : 
Compile directly from source. No third-party Go modules are required.

For Linux :
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o go-relay relay.go

For Windows :
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o go-relay.exe relay.go


## Usage
Syntax: ./go-relay <local_bind_port> <target_ip:target_port>

Example 1: Basic Forward Tunnel

Expose an isolated internal database (10.10.10.50:3306) through a compromised dual-homed web server (192.168.1.100).

Execute on the compromised web server:
./go-relay 9999 10.10.10.50:3306

Result: Tooling pointed at 192.168.1.100:9999 will be seamlessly routed to the internal database.
Example 2: Double Pivoting (Chained Relays)
Route exploit traffic through two compromised systems to reach a deeply segregated target enclave.
1. Hop 2 (Internal Server):
./go-relay 8888 <Final_Target_IP>:445

2. Hop 1 (DMZ Web Server):
./go-relay 9999 <Hop_2_IP>:8888

3. Attacker Execution:
Configure local exploit tooling or proxychains to target <Hop_1_IP>:9999. The TCP stream will automatically traverse Hop 1 and Hop 2 to interface with the final target.

## Disclaimer
This software is developed exclusively for authorized Red Team operations, penetration testing, and security engineering simulations. Users are solely responsible for ensuring they have explicit, written authorization to execute network routing tools within their operational scope.

