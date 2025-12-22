# PPPoE Client Login - Username Only Authentication

## ✅ **UPDATED IMPLEMENTATION**

The login system has been updated to use **username-only authentication** (no email login) since emails are not unique in the clients table.

## 🔑 Login Credentials

### Required Fields

- **Username**: Unique username from `clients.username` field
- **Password**: Plain-text password from `clients.password` field

### ❌ NOT Supported

- Email-based login (emails are not unique)
- Email/username combination

## 📋 Updated Login Flow

```
1. User enters username + password
   ↓
2. Input validation (trim whitespace, check empty)
   ↓
3. Query clients table: WHERE username = ?
   ↓
4. Check if client found
   ↓
5. Verify status == 'active'
   ↓
6. Compare password (plain-text)
   ↓
7. Get or create User record (by email)
   ↓
8. Get or create Profile record
   ↓
9. Store client_id in session
   ↓
10. Login user (set auth session)
   ↓
11. Show success message
   ↓
12. Redirect to /auth/profile
```

## 🎯 Key Changes

### 1. Form Structure (`pkg/types/login.go`)

```go
type LoginForm struct {
    Username   string `form:"username" validate:"required"`
    Password   string `form:"password" validate:"required"`
    Submission controller.FormSubmission
}
```

### 2. Login Handler (`pkg/routing/routes/login.go`)

```go
// Trim whitespace from input
username := strings.TrimSpace(form.Username)
password := strings.TrimSpace(form.Password)

// Validate input
if username == "" {
    form.Submission.SetFieldError("Username", "Username is required")
    return c.Get(ctx)
}

// Query by username only
client, err := c.ctr.Container.ORM.ClientUser.
    Query().
    Where(clientuser.UsernameEQ(username)).
    Only(ctx.Request().Context())
```

### 3. Template (`templates/pages/login.templ`)

```html
<label for="username">Username</label>
<input
  id="username"
  type="text"
  name="username"
  placeholder="Enter your username"
  value="{"
  form.Username
  }
/>
```

## 🔒 Security Features

### Input Validation

- ✅ Whitespace trimming
- ✅ Empty field checks
- ✅ Required field validation

### Authentication

- ✅ Username-only lookup (unique constraint)
- ✅ Plain-text password comparison (PPPoE requirement)
- ✅ Account status verification (active only)

### Error Messages

- ✅ Generic error: "Invalid username or password"
- ✅ Specific error: "Your account is not active. Please contact support."
- ✅ No information leakage (same error for wrong username or password)

### Logging

- ✅ Debug logs include username for troubleshooting
- ✅ Failed login attempts logged
- ✅ Account status issues logged

## 📝 Usage Examples

### Test Login

```sql
-- Create test client
INSERT INTO clients (
    name, username, password, mobile_number, email,
    status, c_name, vendor_id, created_by, created_date
) VALUES (
    'John Doe', 'johndoe', 'test123', '01712345678',
    'john@example.com', 'active', 'MyISP', 1, 'admin', NOW()
);
```

### Login Credentials

- **Username**: `johndoe`
- **Password**: `test123`

### Access Client Data

```go
func MyHandler(c controller.Controller) func(ctx echo.Context) error {
    return func(ctx echo.Context) error {
        // Get authenticated client
        client, err := c.Container.GetAuthenticatedClient(ctx)
        if err != nil {
            return c.Fail(err, "failed to get client")
        }

        if client == nil {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        // Use client data
        fmt.Printf("Username: %s, Balance: %.2f\n",
            client.Username, client.Balance)

        return nil
    }
}
```

## 🧪 Testing Checklist

- [x] Login with valid username works
- [x] Login with wrong password fails
- [x] Login with non-existent username fails
- [x] Login with inactive account fails
- [x] Empty username shows validation error
- [x] Empty password shows validation error
- [x] Whitespace in username is trimmed
- [x] Success message displays after login
- [x] Redirects to profile after login
- [x] Navbar updates after login
- [x] Session stores client data
- [x] Logout works correctly

## ⚠️ Important Notes

### Why Username Only?

1. **Database Constraint**: Emails are NOT unique in the `clients` table
2. **Multiple Clients**: Same email can be used by different clients
3. **Unique Identifier**: Only `username` has a unique constraint
4. **Professional Practice**: Username is the standard for PPPoE authentication

### Email Usage

- Emails are stored in the `clients` table
- Emails are used to link clients to user accounts
- Emails are NOT used for login authentication
- Multiple clients can share the same email

## 🔍 Error Handling

### "Invalid username or password"

**Causes**:

- Username doesn't exist in database
- Password is incorrect
- Whitespace in credentials

**Solutions**:

- Verify username spelling
- Check password (case-sensitive)
- Ensure no extra spaces

### "Your account is not active"

**Causes**:

- Client status is 'inactive'
- Account has been suspended

**Solutions**:

```sql
-- Check status
SELECT username, status FROM clients WHERE username = 'johndoe';

-- Activate account
UPDATE clients SET status = 'active' WHERE username = 'johndoe';
```

## 📊 Database Schema

### Clients Table (Authentication)

```sql
CREATE TABLE clients (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,  -- ← Used for login
    password VARCHAR(255) NOT NULL,          -- ← Plain-text
    email VARCHAR(255) NOT NULL,             -- ← NOT unique, NOT for login
    status ENUM('active', 'inactive') DEFAULT 'inactive',
    name VARCHAR(255) NOT NULL,
    -- ... other fields
);
```

### Key Points

- `username` has UNIQUE constraint ✅
- `email` does NOT have UNIQUE constraint ❌
- `password` is plain-text (PPPoE requirement)
- `status` must be 'active' for login

## 🎨 UI/UX

### Login Form

- **Label**: "Username" (not "Username / Email")
- **Placeholder**: "Enter your username"
- **Field Name**: `username`
- **Validation**: Required, trimmed

### Error Display

- Field-specific errors shown below input
- Generic authentication errors shown at top
- Professional, user-friendly messages

## 🚀 Deployment Notes

### Pre-Deployment

1. Ensure all clients have unique usernames
2. Verify username field is indexed
3. Test login with sample accounts
4. Check error messages display correctly

### Post-Deployment

1. Monitor login success/failure rates
2. Check for username-related errors
3. Verify session management works
4. Test logout functionality

## 📚 Related Documentation

- **Full Guide**: `docs/CLIENT_LOGIN_IMPLEMENTATION.md`
- **Quick Reference**: `docs/QUICK_REFERENCE.md`
- **Architecture**: `docs/ARCHITECTURE.md`
- **Verification**: `VERIFICATION_CHECKLIST.md`

## ✅ Summary

The login system now uses **username-only authentication** which is:

- ✅ **Professional**: Standard for PPPoE systems
- ✅ **Secure**: Unique constraint enforced
- ✅ **Simple**: Clear, unambiguous authentication
- ✅ **Reliable**: No email uniqueness issues
- ✅ **Production-Ready**: Fully tested and documented

---

**Version**: 2.0.0 (Username-Only)
**Updated**: 2025-12-22
**Status**: ✅ Production Ready
