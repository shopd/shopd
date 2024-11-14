-- SessionByUserID fetches a single row
-- name: SessionByUserID :one
select user.email, user.username, user.role, user.verified as user_verified, 
user.disabled, session.*
from user join session on session.user_id = user.user_id
where user.user_id = ? limit 1;
