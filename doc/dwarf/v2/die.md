### 5.2.2 DIE

Each debugging information entry is described by **an identifying tag** and contains **a series of attributes**. 
- The tag specifies the class to which an entry belongs;
- The attributes define the specific characteristics of the entry;

The debugging information entries in Dwarf v2 are intended to exist in the **.debug_info** section of an object file.

>fixme What about Dwarf v4? Is it .zdebug_info?

#### 5.2.2.1 Tag

Tag, specifies what the DIE describes, the set of required tag names is listed in following figure.

![img](assets/clip_image001.png)

#### 5.2.2.2 Attribute

Attribute, fill in details of DIE and further describes the entity.

An attribute has a variety of values: constants (such as function name), variables (such as start address for a function), or references to another DIE (such as for the type of functions’ return value).

The permissive values for an attribute belong to one or more classes of attribute value forms. Each form class may be represented in one or more ways. 

For instance, some attribute values consist of a single piece of constant data. “Constant data” is the class of attribute value that those attributes may have. There’re several representations of constant data, however (one, two, four, eight bytes and variable length data). The particular representation for any given instance of an attribute is encoded along with the attribute name as part of of the information that guides the interpretation of a debugging information entry.

The set of required attribute names is listed in following figure.

![img](assets/clip_image002.png)

**Attribute value forms** may belong to one of the following classes:

1. **Address**, refers to some location in the address space of the described program.
2. **Block**, an arbitrary number of uninterpreted bytes of data.
3. **Constant**, one, two, four or eight bytes of uninterpreted data, or data encoded in LEB128.
4. **Flag**, a small constant that indicates the presence or absence of the an attribute.
5. **Reference**, refers to some member of the set of DIEs that describe the program.
6. **String**, a null-terminated sequence of zero or more (non-null) bytes. Strings maybe represented directly in the DIE or as an offset in a separate string table.

#### 5.2.2.3 Form

Briefly, DIE can be classified into 2 forms: 

1. the one to describe **the data and type**
2. the one to describe **the function and executable code**

One DIE can have parent, siblings and children DIEs, dwarf debugging info is constructed as a tree in which each node is a DIE, several DIE combined to describe a entity in programming language (such as a function).

In following sections, types of DIEs will be described before we dive into dwarf further.

