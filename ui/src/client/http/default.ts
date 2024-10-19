import axios, { AxiosInstance } from 'axios';
import { z } from 'zod';
import { Attendance, AttendanceSchema, Session, SessionSchema, User, UserSchema } from './data';
import setToken from './interceptor';

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;

const client: AxiosInstance = axios.create({
  baseURL: `${BASE_URL.replace(/\/$/, '')}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

client.interceptors.response.use(setToken);

const GetUserAttendancesResponseSchema = z.object({
  attendances: z.array(AttendanceSchema),
});

const GetSessionAttendancesResponseSchema = z.object({
  attendances: z.array(AttendanceSchema),
});

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

export const getSession = async (id: string): Promise<Session> => {
  const response = await client.get(`/sessions/${id}`);
  return SessionResponseSchema.parse(response.data);
};

export type ListSessionsResponse = z.infer<typeof ListSessionsResponseSchema>;

export const listSessions = async (offset: number, pageSize: number): Promise<ListSessionsResponse> => {
  const response = await client.get('/sessions', { params: { offset, pageSize } });
  return ListSessionsResponseSchema.parse(response.data);
};

export const getUserAttendances = async (userId: string): Promise<Attendance[]> => {
  const response = await client.get(`/users/${userId}/attendances`);
  return GetUserAttendancesResponseSchema.parse(response.data).attendances;
};

export const getSessionAttendances = async (sessionId: string): Promise<Attendance[]> => {
  const response = await client.get(`/sessions/${sessionId}/attendances`);
  return GetSessionAttendancesResponseSchema.parse(response.data).attendances;
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

const GetUserAuthResponseSchema = z.object({
  user_id: z.string(),
  user_role: z.enum(['', 'unknown', 'super_admin', 'admin', 'member']),
});

export type GetUserAuthResponse = z.infer<typeof GetUserAuthResponseSchema>;

export const getUserAuth = async (): Promise<GetUserAuthResponse> => {
  const response = await client.get('/auth');
  return GetUserAuthResponseSchema.parse(response.data);
};
