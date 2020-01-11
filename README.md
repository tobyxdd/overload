# Overload

![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-2.png)  
  
Overload is a network load balancing proxy. **But don't get it wrong - it's not a load balancer for your web servers. It's a local proxy that you can run on your PC to distributes traffic across multiple network adapters.**  
  
Imagine the following scenario: You have your own broadband in your home, but you also happen to be able to catch the signal of some free public hotspot nearby (or you somehow know your neighbor's Wi-Fi password). These networks may not be as fast, but you really want to take advantage of this extra bandwidth. As you already know, your Windows PC can only utilize one network at a time (*). Even if you have multiple adapters connected to different networks, traffic will only go through one of them by default.  
  
(*): Try [this](https://docs.microsoft.com/en-us/windows-server/networking/technologies/network-load-balancing) instead if you are using Windows Server  
  
Now, this is where Overload comes in handy. It distributes TCP connections to multiple adapters according to the configuration. Then as long as the applications you use support multi-threaded downloading, you can potentially double, triple or even quadruple the speed (if you have that many adapters & sources of internet available, that is 😜).  
  
## Usage  
  
Place a `config.json` in the working directory (an example is provided in the repo):  
  
```json  
{
  "socks5_listen_addr": "127.0.0.1:1080",
  "interfaces": [
    {
      "name": "Ethernet",
      "weight": 1
    },
    {
      "name": "Wi-Fi 2",
      "weight": 1
    }
  ]
}
```  
  
Greater weights indicate higher chances of being used for a connection. Adjust the weights based on the bandwidth to get best results.

## TODO

 - IPv6 support (this one is a bit tricky, because some adapters may support IPv6 and some others may not)