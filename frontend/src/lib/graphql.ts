import { gql } from '@apollo/client';

// Todo Queries
export const GET_TODOS = gql`
  query GetTodos {
    todos(order_by: { created_at: desc }) {
      id
      user_id
      title
      description
      completed
      created_at
      updated_at
    }
  }
`;

export const GET_TODO = gql`
  query GetTodo($id: Int!) {
    todos_by_pk(id: $id) {
      id
      user_id
      title
      description
      completed
      created_at
      updated_at
    }
  }
`;

// Todo Mutations
export const INSERT_TODO = gql`
  mutation InsertTodo($title: String!, $description: String, $completed: Boolean) {
    insert_todos_one(object: {
      title: $title,
      description: $description,
      completed: $completed
    }) {
      id
      user_id
      title
      description
      completed
      created_at
      updated_at
    }
  }
`;

export const UPDATE_TODO = gql`
  mutation UpdateTodo($id: Int!, $title: String, $description: String, $completed: Boolean) {
    update_todos_by_pk(
      pk_columns: { id: $id },
      _set: {
        title: $title,
        description: $description,
        completed: $completed
      }
    ) {
      id
      user_id
      title
      description
      completed
      created_at
      updated_at
    }
  }
`;

export const DELETE_TODO = gql`
  mutation DeleteTodo($id: Int!) {
    delete_todos_by_pk(id: $id) {
      id
    }
  }
`;

// User Queries (for admin)
export const GET_USERS = gql`
  query GetUsers {
    users(order_by: { created_at: desc }) {
      id
      email
      is_admin
      created_at
      updated_at
    }
  }
`;

export const GET_USER = gql`
  query GetUser($id: Int!) {
    users_by_pk(id: $id) {
      id
      email
      is_admin
      created_at
      updated_at
      todos {
        id
        title
        description
        completed
        created_at
        updated_at
      }
    }
  }
`;
