Each debugging information entry is described by **an identifying tag** and contains **a series of attributes**. 
- The tag specifies the class to which an entry belongs;
- The attributes define the specific characteristics of the entry;

The debugging information entries in Dwarf v2/v3 are intended to exist in the **.debug_info** section of an object file.

> If compiler **compresses debugging information**, compressed debugging information will **be stored in section with “z” prefix**, for example, compressed “.debug_info” will be stored in section “.zdebug_info”.
