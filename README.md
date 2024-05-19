# Config Combinator (combi)

Combi is a simple client that merge diferents configurations stored in a source with a local target configuration file to generate a final usable configuration.

## Daemon flow

```
               ┌─────────────┐                                   
               │             │                                   
 ┌────┬────┬───►  sync time  │                                   
 │    │    │   │             │                                   
 │    │    │   └──────┬──────┘                                   
 │    │    │          │                                          
 │    │    │  ┌───────▼───────┐    ┌──────────┐                  
 │    │    │  │               │    │          │  │local file     
 │    │    no │  get  config  ◄────┤  source  ├─►│git repo       
 │    │    │  │               │    │          │  │...            
 │    │    │  └───────┬───────┘    └──────────┘                  
 │    │    │          │                                          
 │    │    │    ┌─────▼─────┐                                    
 │    │    │    │           │                                    
 │    │    └────┤  update?  │                                    
 │    │         │           │                                    
 │    │         └─────┬─────┘                                    
 │    │               │                                          
 │    │              yes                                         
 │    │               │                                          
 │    │      ┌────────▼────────┐                                 
 │    │      │                 │                                 
 │    │      │  decode config  ◄─────────────┐                   
 │    │      │                 │             │                   
 │    │      └────────┬────────┘             │                   
 │    │               │                      │                   
 │    │      ┌────────▼────────┐             │                   
 │    │      │                 │             │                   
 │    │      │  merge configs  │             │                   
 │    │      │                 │       ┌─────┴─────┐  │libconfig 
 │    │      └────────┬────────┘       │           │  │nginx conf
 │    │               │                │  encoder  ├─►│json      
 │    │    ┌──────────▼──────────┐     │           │  │yaml      
 │    │    │                     │     └─────┬─────┘  │...       
 │    │    │  check  conditions  │           │                   
 │    │    │                     │           │                   
 │    │    └──────────┬──────────┘           │                   
 │    │               │                      │                   
 │    │     ┌─────────┴─────────┐            │                   
 │    │     │                   │            │                   
 │ ┌──┴─────▼────────┐ ┌────────▼─────────┐  │                   
 │ │                 │ │                  │  │                   
 │ │  fail acttions  │ │  encode  config  ◄──┘                   
 │ │                 │ │                  │                      
 │ └─────────────────┘ └────────┬─────────┘                      
 │                              │                                
┌┴───────────────────┐ ┌────────▼─────────┐                      
│                    │ │                  │                      
│  success acttions  ◄─┤  update  config  │                      
│                    │ │                  │                      
└────────────────────┘ └──────────────────┘                      
```
