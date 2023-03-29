### Idea
The idea make this module is config should be vary between deploys.
No grouping after spesific deployment since this is not scale cleanly, an example:
- env.dev.json
- env.prod.json
- env.david.dev.json, etc

It will make harder and harder to managed all those configs. 
Compared to that, configs should be dependent with its deployment.
When deployment does not need anymore, its configs died within it.


### Dependencies
Using [viper](https://github.com/spf13/viper), as it is already has many useful features.
