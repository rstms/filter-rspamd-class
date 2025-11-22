# filter-rspamd-class

Use this filter in an OpenSMTP filter chain after filter-rspamd to apply an
`X-Spam-Class` keyword header to each message with a `X-Spam-Score` header.
This produces an easily matched string field for use by downstream or
client-side filter software.

## Operation
The `X-Spam-Class` header value is based on class names and threshold levels from the configuration.
The score value from the rspamd-generated `X-Spam-Score` is compared with the configuration classes to set the class value.

Default class threshold levels:

Name	    | Threshold	| Condition 
----------- | --------- | -----------------------------------------------------
ham	    | 0		| spam_score < HAM_THRESHOLD
possible    | 3		| HAM_THRESHOLD <= spam_score < POSSIBLE_THRESHOLD
probable    | 10	| POSSIBLE_THRESHOLD <= spam_score < PROBABLE_THRESHOLD
spam	    | 999	| spam_score >= PROBABLE_THRESHOLD


## Configuration
The number of classes, class names, and thresholds are configurable for each recipient email address.

If a config file does not exist, the defaults are used.
If no match is found, the default class values will be used.

When a config file is present, the first matching RCPT-TO address for a message
is used to lookup the class values for that message.


Configuration filename: `/etc/mail/filter_rspamd_classes.json`

### JSON config file example
```
{
    "username@example.org": [
	    { "name": "ham", "score": 0 },
	    { "name": "possible", "score": 3 },
	    { "name": "probable", "score": 10 },
	    { "name": "spam", "score": 999 }
    ],
    "othername@example.org": [
	    { "name": "not_spam", "score": 0 },
	    { "name": "suspected_spam", "score": 10 },
	    { "name": "is_spam", "score": 999 }
    ]
}
```
### Config notes:
The program assumes the array of name, score structs is provided in increasing score order.
The final threshold value is internally set to the maximum possible value; 999 is a placeholder.
