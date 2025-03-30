import axios, { AxiosInstance } from 'axios';
import { z } from 'zod';
import { AdminSession, AdminSessionSchema, Session, UserSchema } from './data';
import setToken from './interceptor';

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;

const client: AxiosInstance = axios.create({
  baseURL: `${BASE_URL.replace(/\/$/, '')}/api/admin`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

client.interceptors.response.use(setToken);

const AdminAllActiveUsersResponseSchema = z.array(UserSchema);

export type AdminAllActiveUsersResponse = z.infer<typeof AdminAllActiveUsersResponseSchema>;

export const adminAllActiveUsers = async (): Promise<AdminAllActiveUsersResponse> => {
  const response = await client.get('/users');
  return AdminAllActiveUsersResponseSchema.parse(response.data);
};

const AdminGetSessionResponseSchema = AdminSessionSchema;

export const adminGetSession = async (id: string): Promise<AdminSession> => {
  const response = await client.get(`/sessions/${id}`);
  return AdminGetSessionResponseSchema.parse(response.data);
};

export const adminCreateSession = async (
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

const AdminListSessionsResponseSchema = z
  .object({
    is_end: z.boolean(),
    sessions: z.array(AdminSessionSchema),
    total_count: z.number(),
  })
  .transform((data) => ({
    isEnd: data.is_end,
    sessions: data.sessions,
    totalCount: data.total_count,
  }));

export type AdminListSessionsResponse = z.infer<typeof AdminListSessionsResponseSchema>;

export const adminListSessions = async (offset: number, pageSize: number): Promise<AdminListSessionsResponse> => {
  const response = await client.get('/sessions', { params: { offset, pageSize } });
  return AdminListSessionsResponseSchema.parse(response.data);
};

export const adminDeleteSession = async (id: string): Promise<void> => {
  await client.delete(`/sessions/${id}`);
};

export const adminCreateSessionForm = async (sessionId: string): Promise<string> => {
  const response = await client.post(`/sessions/${sessionId}/attendance-form`);
  return response.data.form_url;
};

export const adminMarkUsersAsPresent = async (sessionId: string, userIds: string[]): Promise<void> => {
  await client.post(`/sessions/${sessionId}/attendance/manual`, { user_ids: userIds });
};

export const adminLateApplyAttendance = async (sessionId: string, userIds: string[]): Promise<void> => {
  await client.post(`/sessions/${sessionId}/attendance/late`, { user_ids: userIds });
};
