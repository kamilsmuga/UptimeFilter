UptimeFilter
============

Uptime aggregator (per hour) filter for Heka

### Build notes
1. Touch cmake/plugin_loader.cmake
2. Edit plugin_loader.cmake and add <br/>
```
add_external_plugin(git https://github.com/kamilsmuga/UptimeFilter 84da1151b7570e43a6ca6d43b9e0023344be6033)
```
<br/>
where ``` 84da1151b7570e43a6ca6d43b9e0023344be6033 ``` is the commit you want to use
