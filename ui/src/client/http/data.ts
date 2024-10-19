import { z } from 'zod';

export const UserSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    generation: z.number(),
    is_active: z.boolean(),
    email: z.string(),
    external_name: z.string(),
  })
  .transform((data) => ({
    id: data.id,
    name: data.name,
    generation: data.generation,
    isActive: data.is_active,
    email: data.email,
    externalName: data.external_name,
  }));

const AttendanceStatus = z.enum(['not_applied_yet', 'applied', 'ignored']);
export const SessionAttendanceAppliedBy = z.enum(['unknown', 'unspecified', 'manual', 'form']);

export const AdminSessionSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    created_by: z.string(),
    google_form_uri: z.string(),
    google_form_id: z.string(),
    created_at: z.string().transform((str) => new Date(str)),
    starts_at: z.string().transform((str) => new Date(str)),
    score: z.number(),
    attendance_status: AttendanceStatus,
    attendance_applied_by: SessionAttendanceAppliedBy,
  })
  .transform((data) => ({
    ...data,
    createdBy: data.created_by,
    googleFormUri: data.google_form_uri,
    googleFormId: data.google_form_id,
    createdAt: data.created_at,
    startsAt: data.starts_at,
    attendanceStatus: data.attendance_status,
    attendanceAppliedBy: data.attendance_applied_by,
  }));

export const SessionSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    created_by: z.string(),
    created_at: z.string().transform((str) => new Date(str)),
    starts_at: z.string().transform((str) => new Date(str)),
    score: z.number(),
  })
  .transform((data) => ({
    ...data,
    createdBy: data.created_by,
    createdAt: data.created_at,
    startsAt: data.starts_at,
  }));

export const AttendanceSchema = z
  .object({
    id: z.string(),
    session_id: z.string(),
    session_name: z.string(),
    session_score: z.number(),
    session_started_at: z.string().transform((str) => new Date(str)),
    user_id: z.string(),
    user_external_name: z.string(),
    user_generation: z.number(),
    user_joined_at: z.string().transform((str) => new Date(str)),
    created_at: z.string().transform((str) => new Date(str)),
  })
  .transform((data) => ({
    id: data.id,
    sessionId: data.session_id,
    sessionName: data.session_name,
    sessionScore: data.session_score,
    sessionStartedAt: data.session_started_at,
    userId: data.user_id,
    userExternalName: data.user_external_name,
    userGeneration: data.user_generation,
    userJoinedAt: data.user_joined_at,
    createdAt: data.created_at,
  }));

export type User = z.infer<typeof UserSchema>;
export type Session = z.infer<typeof SessionSchema>;
export type AdminSession = z.infer<typeof AdminSessionSchema>;
export type Attendance = z.infer<typeof AttendanceSchema>;
