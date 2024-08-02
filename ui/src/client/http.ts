import axios, { AxiosInstance } from 'axios';
import { z } from 'zod';

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;

const client: AxiosInstance = axios.create({
  baseURL: `${BASE_URL.replace(/\/$/, '')}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

const UserSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    university: z.string(),
    phone: z.string(),
    generation: z.number(),
    is_active: z.boolean(),
  })
  .transform((data) => ({
    id: data.id,
    name: data.name,
    university: data.university,
    phone: data.phone,
    generation: data.generation,
    isActive: data.is_active,
  }));

const SessionSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    hosted_by: z.string(),
    created_by: z.string(),
    google_form_uri: z.string(),
    joinning_users: z.array(z.string()),
    created_at: z.string().transform((str) => new Date(str)),
    starts_at: z.string().transform((str) => new Date(str)),
    score: z.number(),
    is_closed: z.boolean(),
  })
  .transform((data) => ({
    ...data,
    hostedBy: data.hosted_by,
    createdBy: data.created_by,
    googleFormUri: data.google_form_uri,
    joinningUsers: data.joinning_users,
    createdAt: data.created_at,
    startsAt: data.starts_at,
    isClosed: data.is_closed,
  }));

export type User = z.infer<typeof UserSchema>;
export type Session = z.infer<typeof SessionSchema>;

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

export const createUser = async (
  name: string,
  university: string,
  phone: string,
  generation: string,
  isActive: boolean,
): Promise<User> => {
  const response = await client.post('/users', {
    name,
    university,
    phone,
    generation,
    is_active: isActive,
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

export const createSessionForm = async (id: string): Promise<string> => {
  const response = await client.post(`/sessions/${id}/attendance-form`);
  return response.data.form_url;
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
