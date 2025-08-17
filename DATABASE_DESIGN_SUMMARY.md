# IQ Theory - Database Design & Deployment Strategy

## Summary

Based on your IQ Theory app requirements, I've designed a comprehensive PostgreSQL database schema and AWS deployment architecture. Here's what I've created for you:

## Database Choice: PostgreSQL (SQL)

**Why PostgreSQL over NoSQL:**

- **Complex Relationships**: Users, groups, friendships, quiz sessions, and leaderboards have intricate relationships
- **ACID Compliance**: Quiz scores and leaderboards require strong consistency
- **Complex Queries**: Leaderboard calculations with filters (global, groups, friends) are better suited for SQL
- **Data Integrity**: Structured data with clear constraints and validation
- **Analytics**: Better support for reporting and performance analytics

## Key Database Features

### 1. **User Management**

- User accounts with authentication
- Email verification and password hashing
- Profile management with avatars

### 2. **Group System (Classroom Support)**

- Teachers can create groups with join codes
- Students can join groups using codes
- Role-based access (admin/member)
- Group membership tracking

### 3. **Friendship System**

- Friend requests and acceptance
- Status tracking (pending, accepted, declined, blocked)
- Prevents duplicate relationships

### 4. **Quiz System**

- **Normalized parameter storage** (clef types, durations, ledger line options)
- **Dynamic quiz configuration generation** (48 combinations from 11 base records)
- Individual quiz sessions with detailed tracking
- Question-by-question answer recording
- Time tracking for each question and overall quiz
- **Efficient storage** - parameters stored directly in quiz sessions

### 5. **Leaderboard System**

- Materialized views for performance
- Multiple ranking criteria (score, accuracy, speed)
- Global, group, and friend-based leaderboards
- Historical performance tracking
- **Parameter-based filtering** for specific quiz types

## AWS Deployment Architecture

### Core Services

- **ECS Fargate**: Serverless container orchestration
- **RDS PostgreSQL**: Managed database with Multi-AZ
- **Application Load Balancer**: SSL termination and routing
- **S3 + CloudFront**: Static asset delivery (note images)
- **ECR**: Docker image registry

### Cost Estimates

- **Development**: ~$50/month
- **Production**: ~$140/month

### Security Features

- VPC with private subnets
- Encrypted storage and transit
- IAM roles and policies
- Secrets management

## Files Created

1. **`database_schema.sql`** - Complete PostgreSQL schema with:

   - All tables and relationships
   - Indexes for performance
   - Materialized views for leaderboards
   - Sample quiz configurations
   - Triggers for automatic leaderboard updates

2. **`AWS_DEPLOYMENT_GUIDE.md`** - Comprehensive deployment guide with:

   - Detailed AWS service breakdown
   - Cost estimates
   - Docker configurations
   - Terraform infrastructure code
   - CI/CD pipeline setup
   - Security best practices

3. **`server/internal/models/models.go`** - Updated Go models with:
   - All database entities as Go structs
   - JSON and database tags
   - Request/response DTOs
   - Validation tags for API endpoints

## Key Schema Tables

### Core Entities

- `users` - User accounts and profiles
- `groups` - Classroom/study groups
- `group_memberships` - User-group relationships
- `friendships` - Friend relationships
- `clef_types` - Musical clef types (4 records)
- `duration_options` - Quiz duration options (3 records)
- `ledger_line_options` - Ledger line limit options (4 records)
- `quiz_sessions` - Individual quiz attempts (stores parameters directly)
- `quiz_answers` - Question-level tracking
- `leaderboards` - Materialized view for rankings
- `available_quiz_configurations` - Dynamic view generating all combinations

### Quiz Parameter Structure

```sql
-- Normalized parameter storage (11 total records instead of 48)
clef_types: treble, bass, alto, tenor (4 records)
duration_options: 30s, 60s, 120s (3 records)
ledger_line_options: 0, 1, 2, 3 ledger lines (4 records)

-- Dynamic combinations: 4 × 3 × 4 = 48 possible quiz configurations
```

### Example Quiz Configurations

```sql
-- Generated dynamically from parameters
Treble Clef - 30 seconds - No ledger lines
Treble Clef - 60 seconds - Up to 1 ledger line
Bass Clef - 120 seconds - Up to 2 ledger lines
Alto/Tenor variations (all combinations available)
```

## Schema Efficiency Improvements

### **78% Storage Reduction**

- **Before**: 48 quiz configuration records with redundant data
- **After**: 11 normalized parameter records (4 + 3 + 4)
- **Benefit**: Significantly reduced storage and maintenance overhead

### **Dynamic Configuration Generation**

- Quiz configurations generated on-demand from parameter combinations
- Easy to add new options (1 record instead of 12-16)
- Consistent naming and validation through normalization

### **API Simplification**

```go
// Before: Reference pre-stored configuration ID
StartQuizRequest{QuizConfigurationID: uuid}

// After: Direct parameter specification
StartQuizRequest{
    Clef: "treble",
    DurationSeconds: 60,
    MaxLedgerLines: 2
}
```

## Leaderboard Features

### Ranking Criteria

1. **Primary**: Highest score
2. **Secondary**: Fastest completion time
3. **Additional metrics**: Accuracy percentage, total attempts

### Leaderboard Scopes

- **Global**: All users across the platform
- **Groups**: Users within the same classroom/group
- **Friends**: Only among connected friends

## Performance Optimizations

1. **Materialized Views**: Pre-computed leaderboard data
2. **Strategic Indexing**: Optimized for parameter-based queries
3. **Auto-refresh Triggers**: Keep leaderboards current
4. **Normalized Storage**: Faster lookups with smaller parameter tables
5. **Dynamic Views**: `available_quiz_configurations` generates combinations efficiently
6. **Parameter-based Filtering**: Direct filtering on clef, duration, and ledger lines

## Next Steps

1. **Database Setup**: Run the schema SQL on your PostgreSQL instance
2. **Environment Configuration**: Set up AWS services using the deployment guide
3. **API Implementation**: Use the Go models to implement your REST endpoints
4. **Frontend Integration**: Update React components to use the new API structure
5. **Testing**: Implement comprehensive tests for all quiz and leaderboard functionality

## Scalability Considerations

- **Database**: Read replicas for leaderboard queries
- **Caching**: Redis for session data and frequent parameter lookups
- **CDN**: CloudFront for global note image delivery
- **Auto-scaling**: ECS services scale based on demand
- **Parameter Efficiency**: Normalized storage scales better with new quiz options

## Summary

This updated architecture uses an **efficient normalized schema** that provides:

✅ **78% reduction** in configuration storage (11 vs 48 records)  
✅ **Simplified maintenance** - add 1 record instead of 12-16  
✅ **Better performance** - optimized parameter-based queries  
✅ **Enhanced flexibility** - runtime configuration management  
✅ **Future-proof design** - easy to extend with new quiz types

The schema supports your current requirements while providing a solid foundation for future enhancements like real-time multiplayer quizzes, advanced analytics, and mobile app support.
