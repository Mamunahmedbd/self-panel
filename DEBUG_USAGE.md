# Debug Usage Statistics

## Check if data exists in radacct table

```sql
-- Check total records for your username
SELECT COUNT(*) FROM radacct WHERE username = 'm-shakil';

-- Check recent sessions
SELECT
    username,
    acctstarttime,
    acctinputoctets,
    acctoutputoctets,
    (COALESCE(acctinputoctets, 0) + COALESCE(acctoutputoctets, 0)) as total_bytes
FROM radacct
WHERE username = 'm-shakil'
ORDER BY acctstarttime DESC
LIMIT 10;

-- Check what username your client is using
SELECT username, status FROM client_users WHERE id = YOUR_CLIENT_ID;
```

## Verify the logged-in username matches radacct

1. Log in to your ISP panel
2. Check the browser console or network tab
3. Look for the username being sent in requests
4. Compare with the `username` field in `radacct` table

## Test the aggregation manually

```sql
-- Today's usage (adjust timezone as needed)
SELECT
    SUM(COALESCE(acctinputoctets, 0) + COALESCE(acctoutputoctets, 0)) as today_bytes
FROM radacct
WHERE username = 'm-shakil'
AND acctstarttime >= CURRENT_DATE;

-- This month's usage
SELECT
    SUM(COALESCE(acctinputoctets, 0) + COALESCE(acctoutputoctets, 0)) as month_bytes
FROM radacct
WHERE username = 'm-shakil'
AND acctstarttime >= DATE_TRUNC('month', CURRENT_DATE);
```

## Expected Results

For the session you provided:

- Download: 4,903 bytes
- Upload: 3,172 bytes
- Total: 8,075 bytes (7.88 KB)

If this session is within today/this week/this month, you should see at least 8 KB in the respective period.
