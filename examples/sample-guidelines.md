# TypeScript Team Guidelines

## Naming Conventions

Use PascalCase for types, interfaces, and classes. Use camelCase for functions, variables, and methods.

### Good Example
```typescript
interface UserProfile {
  firstName: string;
  lastName: string;
}

function getUserProfile(userId: string): UserProfile {
  // implementation
}
```

### Bad Example
```typescript
interface user_profile {
  first_name: string;
  last_name: string;
}

function get_user_profile(user_id: string): user_profile {
  // implementation
}
```

## Type Safety

- Always use explicit types for function parameters and return values
- Avoid using `any` type unless absolutely necessary
- Prefer type assertions over angle bracket syntax

### Rules
- no-explicit-any
- prefer-as-const
- explicit-function-return-type

## Error Handling

All async functions must include proper error handling with try-catch blocks.

### Example
```typescript
async function fetchUserData(id: string): Promise<User> {
  try {
    const response = await fetch(`/api/users/${id}`);
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch user data:', error);
    throw new Error('User data unavailable');
  }
}
```

## Import/Export Practices

- Use named exports instead of default exports
- Group imports by type (external libraries, internal modules, types)
- Use absolute imports when possible

### Good Practice
```typescript
import { Component, OnInit } from '@angular/core';
import { UserService } from '../services/user.service';
import type { User } from '../types/user.types';
```