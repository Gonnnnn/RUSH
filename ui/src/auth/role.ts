import { GetUserAuthResponse } from '../client/http/default';

export enum Role {
  UNKNOWN = 'unknown',
  ADMIN = 'admin',
  MEMBER = 'member',
}

export const convertGetUserAuthRole = (role: GetUserAuthResponse['user_role']): Role => {
  switch (role) {
    case 'unknown':
      return Role.UNKNOWN;
    case 'super_admin':
      return Role.ADMIN;
    case 'admin':
      return Role.ADMIN;
    case 'member':
      return Role.MEMBER;
    default:
      return Role.UNKNOWN;
  }
};
