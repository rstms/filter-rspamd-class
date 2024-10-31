# filter-rspamd-class

## purpose
Use this filter in a filter chain after filter-rspamd to apply a keyword spam class to each message

## operation
Adds a header: 'X-Spam-Class: SPAMCLASS'
based on threshold levels compared with rspamd's X-Spam-Score header

Default class threshold levels:

Name	    | Threshold	| Condition 
----------- | --------- | -----------------------------------------------------
ham	    | 0		| spam_score < HAM_THRESHOLD
possible    | 3		| HAM_THRESHOLD <= spam_score < POSSIBLE_THRESHOLD
probable    | 10	| POSSIBLE_THRESHOLD <= spam_score < PROBABLE_THRESHOLD
spam	    | 999	| spam_score >= PROBABLE_THRESHOLD


## configuration
Class names and thresholds are configurable per recipient email address
If a config file does not exist, the defaults are used.

When a config file is present, the first matching RCPT-TO address for a message
is used to lookup the class values for that messages.
If no match is found, the default class values will be used.

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
