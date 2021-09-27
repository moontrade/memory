libart [![Build Status](https://travis-ci.org/armon/libart.png)](https://travis-ci.org/armon/libart)
=========

This library provides a C99 implementation of the Adaptive Radix
Tree or ART. The ART operates similar to a traditional radix tree but
avoids the wasted space of internal nodes by changing the node size.
It makes use of 4 node sizes (4, 16, 48, 256), and can guarantee that
the overhead is no more than 52 bytes per key, though in practice it is
much lower.

As a radix tree, it provides the following:
* O(k) operations. In many cases, this can be faster than a hash table since
  the hash function is an O(k) operation, and hash tables have very poor cache locality.
* Minimum / Maximum value lookups
* Prefix compression
* Ordered iteration
* Prefix based iteration


References
----------

Related works:

* [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases](http://www-db.in.tum.de/~leis/papers/ART.pdf)
