package tools

const AssetsSchema = `{
    "data": {
        "__schema": {
            "queryType": {
                "name": "Query"
            },
            "types": [
                {
                    "kind": "OBJECT",
                    "name": "Query",
                    "description": "Query.",
                    "fields": [
                       {
                            "name": "assetSearch",
                            "description": "Retrieve a list of assets.",
                            "args": [
                                {
                                    "name": "first",
                                    "description": "Page size.",
                                    "type": {
                                        "kind": "SCALAR",
                                        "name": "Int",
                                        "ofType": null
                                    },
                                    "defaultValue": "20"
                                },
                                {
                                    "name": "after",
                                    "description": "Cursor for the last page item.",
                                    "type": {
                                        "kind": "SCALAR",
                                        "name": "String",
                                        "ofType": null
                                    },
                                    "defaultValue": null
                                },
                                {
                                    "name": "where",
                                    "description": "Filtering.",
                                    "type": {
                                        "kind": "INPUT_OBJECT",
                                        "name": "AssetWhereInput",
                                        "ofType": null
                                    },
                                    "defaultValue": null
                                },
                                {
                                    "name": "inOrganizations",
                                    "description": "Which Organizations to search.",
                                    "type": {
                                        "kind": "LIST",
                                        "name": null,
                                        "ofType": {
                                            "kind": "NON_NULL",
                                            "name": null,
                                            "ofType": {
                                                "kind": "SCALAR",
                                                "name": "ID",
                                                "ofType": null
                                            }
                                        }
                                    },
                                    "defaultValue": null
                                }
                            ],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "OBJECT",
                                    "name": "AssetConnection",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
				},
                {
                    "kind": "INPUT_OBJECT",
                    "name": "AssetWhereInput",
                    "description": "Asset filtering options.\n\nWithin the same input, the fields will be groups together based on their closest AND or OR ancestor.\n\nIf no ancestor exists the operator will default to AND.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "and",
                            "description": "Filter asset by a combination of filters, where asset meets all supplied criteria.",
                            "type": {
                                "kind": "LIST",
                                "name": null,
                                "ofType": {
                                    "kind": "NON_NULL",
                                    "name": null,
                                    "ofType": {
                                        "kind": "INPUT_OBJECT",
                                        "name": "AssetWhereInput",
                                        "ofType": null
                                    }
                                }
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "or",
                            "description": "Filter asset by a combination of filters, where asset meets any of the supplied criteria.",
                            "type": {
                                "kind": "LIST",
                                "name": null,
                                "ofType": {
                                    "kind": "NON_NULL",
                                    "name": null,
                                    "ofType": {
                                        "kind": "INPUT_OBJECT",
                                        "name": "AssetWhereInput",
                                        "ofType": null
                                    }
                                }
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "text",
                            "description": "Filter asset by text.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "TextFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "operatingSystem",
                            "description": "Filter asset by Operating System.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "AssetOperatingSystemWhereInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "name",
                            "description": "Filter asset by Name",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "systemInfo",
                            "description": "Filter asset by System Info",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "AssetSystemInfoWhereInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "description",
                            "description": "Filter asset by Description",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "cpu",
                            "description": "Filter asset by CPU Info.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "AssetCpuWhereInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "externalIpAddress",
                            "description": "Filter asset by external IP address.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "localTimezone",
                            "description": "Filter asset by local time zone.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lastBootedAt",
                            "description": "Filter asset by last boot time.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "DateTimeFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "SCALAR",
                    "name": "BigInt",
                    "description": "Represents integers that are larger than the basic Int scalar.",
                    "fields": null,
                    "inputFields": null,
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "BigIntFilterInput",
                    "description": "Integer property filtering options for values bigger than Int covers. If both gte and lte are populated, a between range will be used.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "equals",
                            "description": "Value is an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notEquals",
                            "description": "Value is not an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gt",
                            "description": "Value is greater than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gte",
                            "description": "Value is greater than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lt",
                            "description": "Value is less than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lte",
                            "description": "Value is less than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "TextFilterInput",
                    "description": "Filter by text across all searchable fields.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "contains",
                            "description": "Value contains a match",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "OBJECT",
                    "name": "AssetConnection",
                    "description": "Asset connection.",
                    "fields": [
                        {
                            "name": "totalCount",
                            "description": "Total asset count.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "SCALAR",
                                    "name": "Int",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "edges",
                            "description": "Asset edges.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "LIST",
                                    "name": null,
                                    "ofType": {
                                        "kind": "NON_NULL",
                                        "name": null,
                                        "ofType": {
                                            "kind": "OBJECT",
                                            "name": "AssetEdge",
                                            "ofType": null
                                        }
                                    }
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "nodes",
                            "description": "Asset nodes.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "LIST",
                                    "name": null,
                                    "ofType": {
                                        "kind": "NON_NULL",
                                        "name": null,
                                        "ofType": {
                                            "kind": "OBJECT",
                                            "name": "Asset",
                                            "ofType": null
                                        }
                                    }
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "AssetOperatingSystemWhereInput",
                    "description": "Filter asset by Operating system.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "name",
                            "description": "Filter asset by Operating system name.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "version",
                            "description": "Filter asset by Operating system version.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "architecture",
                            "description": "Filter asset by Operating system architecture.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "installedOn",
                            "description": "Filter asset by date Operating system was installed on.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "DateTimeFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "type",
                            "description": "Filter asset by Operating system Type.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "AssetOperatingSystemTypeFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "StringFilterInput",
                    "description": "String property filtering options.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "contains",
                            "description": "Value is a partial match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notContains",
                            "description": "Value is not a partial match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "equals",
                            "description": "Value is an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notEquals",
                            "description": "Value is not an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "startsWith",
                            "description": "Value is begining match",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "endsWith",
                            "description": "Value is end of match",
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "DateTimeFilterInput",
                    "description": "DateTime property filtering options. If both gte and lte are populated, a between range will be used.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "equals",
                            "description": "Value is an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gt",
                            "description": "Value is greater than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gte",
                            "description": "Value is greater than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lt",
                            "description": "Value is less than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lte",
                            "description": "Value is less than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notEquals",
                            "description": "Value is not an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "AssetOperatingSystemTypeFilterInput",
                    "description": "Filter enum for operating system type.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "equals",
                            "description": "Value is exact match.",
                            "type": {
                                "kind": "ENUM",
                                "name": "AssetOperatingSystemTypeInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notEquals",
                            "description": "Value is not an exact match.",
                            "type": {
                                "kind": "ENUM",
                                "name": "AssetOperatingSystemTypeInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "ENUM",
                    "name": "AssetOperatingSystemTypeInput",
                    "description": "Enum values for asset operating system type.",
                    "fields": null,
                    "inputFields": null,
                    "interfaces": null,
                    "enumValues": [
                        {
                            "name": "DARWIN",
                            "description": "Darwin.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "LINUX",
                            "description": "Linux.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "WINDOWS",
                            "description": "Windows.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "UNKNOWN",
                            "description": "Unknown.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "AssetSystemInfoWhereInput",
                    "description": "Filter asset by System Info.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "cpuCores",
                            "description": "Filter asset by number of system CPU cores",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "IntFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "cpuName",
                            "description": "Filter asset by system CPU Name",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "hostname",
                            "description": "Filter asset by system Hostname",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "netBiosName",
                            "description": "Filter asset by system Net BIOS Name",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "manufacturer",
                            "description": "Filter asset by system Manufacturer",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "model",
                            "description": "Filter asset by system Model",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "serialNumber",
                            "description": "Filter asset by system SerialNumber",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "IntFilterInput",
                    "description": "Integer property filtering options. If both gte and lte are populated, a between range will be used.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "equals",
                            "description": "Value is an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "notEquals",
                            "description": "Value is not an exact match.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gt",
                            "description": "Value is greater than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "gte",
                            "description": "Value is greater than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lt",
                            "description": "Value is less than input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "lte",
                            "description": "Value is less than or equal to input.",
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "INPUT_OBJECT",
                    "name": "AssetCpuWhereInput",
                    "description": "Filter asset by CPU Info.",
                    "fields": null,
                    "inputFields": [
                        {
                            "name": "maxClockSpeed",
                            "description": "Filter asset by maximum CPU clock speed.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "IntFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "cores",
                            "description": "Filter asset by number of system CPU cores.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "IntFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "count",
                            "description": "Filter asset by CPU count.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "IntFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "name",
                            "description": "Filter asset by system CPU Name.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        },
                        {
                            "name": "type",
                            "description": "Filter asset by CPU type.",
                            "type": {
                                "kind": "INPUT_OBJECT",
                                "name": "StringFilterInput",
                                "ofType": null
                            },
                            "defaultValue": null
                        }
                    ],
                    "interfaces": null,
                    "enumValues": null,
                    "possibleTypes": null
                },
                {
                    "kind": "OBJECT",
                    "name": "AssetEdge",
                    "description": "Asset edge.",
                    "fields": [
                        {
                            "name": "node",
                            "description": "Asset node.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "OBJECT",
                                    "name": "Asset",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "cursor",
                            "description": "Asset cursor.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "SCALAR",
                                    "name": "String",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },  
                {
                    "kind": "OBJECT",
                    "name": "Asset",
                    "description": "Asset entity.",
                    "fields": [
                        {
                            "name": "id",
                            "description": "Asset unique identifier.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "SCALAR",
                                    "name": "ID",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "name",
                            "description": "Asset name.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "description",
                            "description": "Asset description.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "operatingSystemInfo",
                            "description": "Asset operating system.",
                            "args": [],
                            "type": {
                                "kind": "OBJECT",
                                "name": "AssetOperatingSystem",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "lastBootedAt",
                            "description": "Asset last boot time.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "systemInfo",
                            "description": "System Info.",
                            "args": [],
                            "type": {
                                "kind": "OBJECT",
                                "name": "AssetSystemInfo",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "cpu",
                            "description": "Information for CPUs within the asset.",
                            "args": [],
                            "type": {
                                "kind": "OBJECT",
                                "name": "AssetCpus",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "localTimezone",
                            "description": "Asset local timezone, for example 'Eastern Standard Time'.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "operatingSystem",
                            "description": "Asset operating system.",
                            "args": [],
                            "type": {
                                "kind": "OBJECT",
                                "name": "OperatingSystem",
                                "ofType": null
                            },
                            "isDeprecated": true,
                            "deprecationReason": "Replaced by AssetOperatingSystem."
                        },
                        {
                            "name": "externalIpAddress",
                            "description": "Asset external IP address.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
					],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },
	            {
                    "kind": "OBJECT",
                    "name": "AssetOperatingSystem",
                    "description": "Asset Operating system.",
                    "fields": [
                        {
                            "name": "name",
                            "description": "Operating system name.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "version",
                            "description": "Operating system version.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "installedOn",
                            "description": "Date operating system was installed.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "DateTime",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "type",
                            "description": "Operating system type.",
                            "args": [],
                            "type": {
                                "kind": "ENUM",
                                "name": "AssetOperatingSystemType",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "architecture",
                            "description": "Operating system architecture.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                }, 
                {
                    "kind": "ENUM",
                    "name": "AssetOperatingSystemType",
                    "description": "Operating system type.",
                    "fields": null,
                    "inputFields": null,
                    "interfaces": null,
                    "enumValues": [
                        {
                            "name": "DARWIN",
                            "description": "Darwin.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "LINUX",
                            "description": "Linux.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "WINDOWS",
                            "description": "Windows.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "UNKNOWN",
                            "description": "Unknown.",
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "possibleTypes": null
                }, 
                {
                    "kind": "OBJECT",
                    "name": "AssetSystemInfo",
                    "description": "System Information related to the asset.",
                    "fields": [
                        {
                            "name": "manufacturer",
                            "description": "Hardware vendor.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "model",
                            "description": "Hardware model.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "serialNumber",
                            "description": "Device serial number.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "cpuCores",
                            "description": "Number of logical CPU cores available to the system.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": true,
                            "deprecationReason": "Moved to AssetCpu.cpuCores"
                        },
                        {
                            "name": "cpuName",
                            "description": "CPU brand string, contains vendor and model.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": true,
                            "deprecationReason": "Moved to AssetCpus.cpuName"
                        },
                        {
                            "name": "memoryTotalSizeBytes",
                            "description": "Total physical memory in bytes.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "BigInt",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "memoryTotalSizeGB",
                            "description": "Total physical memory in GB.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Float",
                                "ofType": null
                            },
                            "isDeprecated": true,
                            "deprecationReason": "Use memoryTotalSizeBytes instead"
                        },
                        {
                            "name": "hostname",
                            "description": "Network hostname including domain.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "netBiosName",
                            "description": "Friendly computer name.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                }, 
                {
                    "kind": "OBJECT",
                    "name": "AssetCpus",
                    "description": "CPU Information related to the asset.",
                    "fields": [
                        {
                            "name": "cpus",
                            "description": "Detailed list of CPUs.",
                            "args": [],
                            "type": {
                                "kind": "LIST",
                                "name": null,
                                "ofType": {
                                    "kind": "OBJECT",
                                    "name": "AssetCpu",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "maxClockSpeed",
                            "description": "Maximum cpu frequency.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "count",
                            "description": "Number of physical and logical processors detected on the asset.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "cores",
                            "description": "Number of logical CPU cores available on the asset.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "name",
                            "description": "CPU brand string, contains vendor and model.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "type",
                            "description": "The name of your asset's processor and its speed.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },   
                {
                    "kind": "OBJECT",
                    "name": "AssetCpu",
                    "description": "Asset CPU information.",
                    "fields": [
                        {
                            "name": "cores",
                            "description": "Number of logical CPU cores available on the asset.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "cpuId",
                            "description": "The ID of the CPU.",
                            "args": [],
                            "type": {
                                "kind": "NON_NULL",
                                "name": null,
                                "ofType": {
                                    "kind": "SCALAR",
                                    "name": "ID",
                                    "ofType": null
                                }
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "maxClockSpeedMhz",
                            "description": "Maximum CPU frequency in megahertz.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "Int",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "model",
                            "description": "The model of the CPU.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "name",
                            "description": "CPU socket designation.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "type",
                            "description": "The processor type, such as Central, Math, or Video.",
                            "args": [],
                            "type": {
                                "kind": "ENUM",
                                "name": "AssetCpuType",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },  
                {
                    "kind": "OBJECT",
                    "name": "OperatingSystem",
                    "description": "Operating system.",
                    "fields": [
                        {
                            "name": "name",
                            "description": "Operating system name.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "version",
                            "description": "Operating system version.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        },
                        {
                            "name": "architecture",
                            "description": "Operating system architecture.",
                            "args": [],
                            "type": {
                                "kind": "SCALAR",
                                "name": "String",
                                "ofType": null
                            },
                            "isDeprecated": false,
                            "deprecationReason": null
                        }
                    ],
                    "inputFields": null,
                    "interfaces": [],
                    "enumValues": null,
                    "possibleTypes": null
                },              
				]
			}
		}
}`
