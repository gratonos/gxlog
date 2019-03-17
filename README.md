# gxlog #

This gxlog is a simplified version of [gxlog](https://github.com/gxlog/gxlog).

```
+----------------------------------------------------------------+
|                             Logger                             |
|                          Level [Filter]                        |
| Record                                                         |
|   | +-------+------------------------------------------------+ |
|   |-| Slot0 | Formatter Writer Level [Filter] [ErrorHandler] | |
|   | +-------+------------------------------------------------+ |
|   |-|  ...  |              ...                               | |
|   | +-------+------------------------------------------------+ |
|   \-| Slot7 | Formatter Writer Level [Filter] [ErrorHandler] | |
|     +-------+------------------------------------------------+ |
+----------------------------------------------------------------+
```
