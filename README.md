# PackGen
Go based Packet Generator


# Help
bash# ./PackGen --help

Usage of ./PackGen:

  -c int
  
    	Number of connections to be generated with auto-generated MAC and IP (default 10)
      
  -di string
  
    	destination ip (default "127.0.0.1")
      
  -dm string
  
    	destination MAC address (default "FF:FF:FF:FF:FF:FF")
      
  -i string
  
    	device name (default "lo")
      
  -r int
  
    	Rate of packet generation needed in seconds (default 30)
      
  -si string
  
    	source ip (default "20.0.0.1")
      
  -sm string
  
    	Src MAC address (default "00:01:02:00:00:00")
      
      
      
# Example Usage
bash# ./PackGen -c 1000 -di 40.40.40.1 -dm 00:0c:29:f2:cf:73 -i ens224 -r 40
CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

CurrentRate is 40

Done generating packets.

Good Bye!!!

bash#

