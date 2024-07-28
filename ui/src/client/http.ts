import axios, { AxiosInstance } from 'axios';
import { z } from 'zod';

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;

// set cors
const client: AxiosInstance = axios.create({
  baseURL: `${BASE_URL.replace(/\/$/, '')}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
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

const AttendanceSchema = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    session_ids: z.array(z.string()),
    created_at: z.string().transform((str) => new Date(str)),
    created_by: z.string(),
  })
  .transform((data) => ({
    ...data,
    sessionIds: data.session_ids,
    createdAt: data.created_at,
    createdBy: data.created_by,
  }));

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

const AttendancesResponseSchema = z.array(AttendanceSchema);

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

export const createSession = async (
  name: string,
  description: string,
  startsAt: Date,
  score: number,
): Promise<Session> => {
  const response = await client.post('/sessions', {
    name,
    description,
    starts_at: {
      year: startsAt.getFullYear(),
      month: startsAt.getMonth() + 1,
      day: startsAt.getDate(),
      hour: startsAt.getHours(),
      minute: startsAt.getMinutes(),
    },
    score,
  });
  return response.data.id;
};

export const createSessionForm = async (id: string): Promise<string> => {
  const response = await client.post(`/sessions/${id}/attendance-form`);
  return response.data.form_url;
};

export const getAttendances = async (): Promise<Attendance[]> => {
  const response = await client.get('/attendances');
  return AttendancesResponseSchema.parse(response.data);
};
