# dic
Dictory under command line

This is a dictory under command line,you can special remote url include traslate and sound.For Example:config.json


Format
==================
<pre>
config.json format followed:  
{  
    "{scheme}"://the second args  
    {  
        "url":"remote url",//remote url  
        "filter":[// how to find translation or sound  
          ".DEF",  
          "index:0"  
        ],  
        "enable":true,  
        "dir":["/Users/jeekyuan/english/dic/{word}"]//cache dir,one of array must readed and writed  
    }  
}  
</pre>

<pre>
dic {scheme} {word}  
the scheme is "default" if you type "dic {word}"  
</pre>

Operate
=============
You can set cmd in config.json  
cmd followed:  
<pre>
"player" : use dic's player to play sound  
"text"   : do nothing  
"sh"     : running a command. for example: sh mplayer.exe {path}  
</pre>

Variable
==============
<pre>
{word}   : word searched  
{path}   : path cached  
{scheme} : config.json's key  
{suff}   : the suffer of found result by Filter  
</pre>


Play Sournd
==============
Myself decoder onlay decode MPEG-1, please use system's player if you want to play sound  
the Demo(config.json) is running on Mac  
