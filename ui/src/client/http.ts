import axios, { AxiosInstance } from 'axios';
import { z } from 'zod';
import { replaceCookieHeader, setAuthToken } from './auth';

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;

const client: AxiosInstance = axios.create({
  baseURL: `${BASE_URL.replace(/\/$/, '')}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

client.interceptors.response.use((response) => {
  const newToken = response.headers[replaceCookieHeader];
  if (newToken) {
    setAuthToken(newToken);
  }
  return response;
});

const UserSchema = z
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

const AttendanceStatusSchema = z.enum(['not_applied_yet', 'applied', 'ignored']);

const SessionSchema = z
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
    attendance_status: AttendanceStatusSchema,
  })
  .transform((data) => ({
    ...data,
    createdBy: data.created_by,
    googleFormUri: data.google_form_uri,
    googleFormId: data.google_form_id,
    createdAt: data.created_at,
    startsAt: data.starts_at,
    attendanceStatus: data.attendance_status,
  }));

const AttendanceSchema = z
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

const GetUserAttendancesResponseSchema = z.object({
  attendances: z.array(AttendanceSchema),
});

const GetSessionAttendancesResponseSchema = z.object({
  attendances: z.array(AttendanceSchema),
});

export type User = z.infer<typeof UserSchema>;
export type Session = z.infer<typeof SessionSchema>;
export type Attendance = z.infer<typeof AttendanceSchema>;

const ListUsersResponseSchema = z
  .object({
    is_end: z.boolean(),
    users: z.array(UserSchema),
    total_count: z.number(),
  })
  .transform((data) => ({
    isEnd: data.is_end,
    users: data.users,
    totalCount: data.total_count,
  }));
const UserResponseSchema = UserSchema;

const SessionsResponseSchema = z.array(SessionSchema);
const SessionResponseSchema = SessionSchema;

const ListSessionsResponseSchema = z
  .object({
    is_end: z.boolean(),
    sessions: z.array(SessionSchema),
    total_count: z.number(),
  })
  .transform((data) => ({
    isEnd: data.is_end,
    sessions: data.sessions,
    totalCount: data.total_count,
  }));

export type ListUsersReponse = z.infer<typeof ListUsersResponseSchema>;

export const listUsers = async (offset: number, pageSize: number): Promise<ListUsersReponse> => {
  const response = await client.get('/users', { params: { offset, pageSize } });
  return ListUsersResponseSchema.parse(response.data);
};

export const getUser = async (id: string): Promise<User> => {
  const response = await client.get(`/users/${id}`);
  return UserResponseSchema.parse(response.data);
};

export const createUser = async (name: string, generation: string, isActive: boolean, email: string): Promise<User> => {
  const response = await client.post('/users', {
    name,
    generation,
    is_active: isActive,
    email,
  });
  return UserResponseSchema.parse(response.data);
};

export const getSession = async (id: string): Promise<Session> => {
  const response = await client.get(`/sessions/${id}`);
  return SessionResponseSchema.parse(response.data);
};

export const getSessions = async (): Promise<Session[]> => {
  const response = await client.get('/sessions');
  return SessionsResponseSchema.parse(response.data);
};

export type ListSessionsResponse = z.infer<typeof ListSessionsResponseSchema>;

export const listSessions = async (offset: number, pageSize: number): Promise<ListSessionsResponse> => {
  const response = await client.get('/sessions', { params: { offset, pageSize } });
  return ListSessionsResponseSchema.parse(response.data);
};

export const createSession = async (
  name: string,
  description: string,
  startsAt: Date,
  score: number,
): Promise<Session> => {
  const response = await client.post('/sessions', {
    name,
    description,
    starts_at: startsAt.toISOString(),
    score,
  });
  return response.data.id;
};

export const deleteSession = async (id: string): Promise<void> => {
  await client.delete(`/sessions/${id}`);
};

export const createSessionForm = async (sessionId: string): Promise<string> => {
  const response = await client.post(`/sessions/${sessionId}/attendance-form`);
  return response.data.form_url;
};

export const closeSession = async (sessionId: string): Promise<void> => {
  await client.post(`/sessions/${sessionId}/attendance`);
};

export const getUserAttendances = async (userId: string): Promise<Attendance[]> => {
  const response = await client.get(`/users/${userId}/attendances`);
  return GetUserAttendancesResponseSchema.parse(response.data).attendances;
};

export const getSessionAttendances = async (sessionId: string): Promise<Attendance[]> => {
  const response = await client.get(`/sessions/${sessionId}/attendances`);
  return GetSessionAttendancesResponseSchema.parse(response.data).attendances;
};

export const markUsersAsPresent = async (sessionId: string, userIds: string[]): Promise<void> => {
  await client.post(`/sessions/${sessionId}/present`, { user_ids: userIds });
};

const GetHalfYearAttendancesResponseSchema = z
  .object({
    sessions: z.array(
      z.object({
        id: z.string(),
        name: z.string(),
        started_at: z.string().transform((str) => new Date(str)),
      }),
    ),
    users: z.array(
      z.object({
        id: z.string(),
        name: z.string(),
        generation: z.number(),
      }),
    ),
    attendances: z.array(AttendanceSchema),
  })
  .transform((data) => ({
    sessions: data.sessions.map((session) => ({
      id: session.id,
      name: session.name,
      startedAt: session.started_at,
    })),
    users: data.users,
    attendances: data.attendances,
  }));

export type GetHalfYearAttendancesResponse = z.infer<typeof GetHalfYearAttendancesResponseSchema>;

export const getHalfYearAttendances = async (): Promise<GetHalfYearAttendancesResponse> => {
  const response = await client.get('/attendances/half-year');
  return GetHalfYearAttendancesResponseSchema.parse(response.data);
};

export const signIn = async (token: string): Promise<string> => {
  const response = await client.post('/sign-in', { token });
  if (response.status === 200) {
    return response.data.token;
  }
  throw new Error('Failed to sign in');
};

export const checkAuth = async (): Promise<boolean> => {
  const response = await client.get('/auth');
  return response.status === 200;
};

export const getUserId = async (): Promise<string> => {
  const response = await client.get('/auth');
  return response.data.user_id;
};
