import React, { useEffect, useState } from'react';
import { API, showError, showSuccess } from '../helpers';
import { Route, Routes, useNavigate, Navigate } from'react-router-dom';
import { formatTime } from './utils';
import {
  Button,
  Form,
  Popconfirm,
  Space,
  Table,
  Tag,
  Tooltip,
} from '@douyinfe/semi-ui';
import { IconSearch, IconSidebar, IconChevronDown } from '@douyinfe/semi-icons';
import { ITEMS_PER_PAGE } from '../constants';
import { useTranslation } from'react-i18next';

// 提取渲染类型的函数到外部，提高代码可读性
const renderType = (type, t) => {
  switch (type) {
    case 'text_chat':
      return <Tag size='large'>{t('对话')}</Tag>;
    case 'image_generate':
      return (
        <Tag color='yellow' size='large'>
          {t('image_generate')}
        </Tag>
      );
    default:
      return (
        <Tag color='red' size='large'>
          {t('未分类')}
        </Tag>
      );
  }
};

// 提取渲染角色的函数到外部，提高代码可读性
const renderRole = (role, t) => {
  switch (role) {
    case 1:
      return <Tag size='large'>{t('普通用户')}</Tag>;
    case 10:
      return (
        <Tag color='yellow' size='large'>
          {t('管理员')}
        </Tag>
      );
    case 100:
      return (
        <Tag color='orange' size='large'>
          {t('超级管理员')}
        </Tag>
      );
    default:
      return (
        <Tag color='red' size='large'>
          {t('未知身份')}
        </Tag>
      );
  }
};

const handleAsyncError = (error) => {
  showError(error.message);
};

const TableConversation = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const columns = [
    {
      title: t('用户名'),
      dataIndex: 'username',
    },
    {
      title: t('会话标题'),
      dataIndex: 'title',
    },
    {
      title: t('会话类型'),
      dataIndex: 'type',
      render: (text, record, index) => {
        return <div>{renderType(text, t)}</div>;
      },
    },
    {
      title: t('创建时间'),
      dataIndex: 'created_at',
      render: (text, record, index) => {
        return <div>{formatTime(new Date(text))}</div>;
      },
    },
    {
      title: t('修改时间'),
      dataIndex: 'updated_at',
      render: (text, record, index) => {
        return <div>{formatTime(new Date(text))}</div>;
      },
    },
    {
      title: '',
      dataIndex: 'id',
      render: (text,record) => (
        <Button theme="borderless" onClick={() => {
          navigate("/conversationdetail", {
            state: { id: text,type:record.type }
          })
        }}>详情</Button>
      ),
    },
  ];

  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activePage, setActivePage] = useState(1);
  const [pageSize, setPageSize] = useState(ITEMS_PER_PAGE);
  const [searchTitle, setSearchTitle] = useState('');
  const [searchUsername, setSearchUsername] = useState('');

  const [searching, setSearching] = useState(false);
  const [searchGroup, setSearchGroup] = useState('');
  const [groupOptions, setGroupOptions] = useState([]);
  const [userCount, setUserCount] = useState(ITEMS_PER_PAGE);
  const [showAddUser, setShowAddUser] = useState(false);
  const [showEditUser, setShowEditUser] = useState(false);
  const [editingUser, setEditingUser] = useState({ id: undefined });

  // 提取删除记录的函数到外部，提高代码可读性
  const removeRecord = (key) => {
    setUsers(prevUsers => {
      return prevUsers.map(user => {
        if (user.id === key) {
          return { ...user, DeletedAt: new Date() };
        }
        return user;
      });
    });
  };

  // 提取设置用户格式的函数到外部，提高代码可读性
  const setUserFormat = (users) => {
    const newUsers = users.map(user => ({ ...user, key: user.id }));
    setUsers(newUsers);
  };

  const loadUsers = async (startIdx, pageSize) => {
    console.log("pageL",activePage)

    try {
      const res = await API.get(
        `/api/conversation/list_all?page=${startIdx}&page_size=${pageSize}`,
      );
      const { success, message, data } = res.data;
      if (success) {
        const newPageData = data.items;
        setActivePage(data.page);
        setUserCount(data.total);
        setUserFormat(newPageData);

      } else {
        showError(message);
      }
    } catch (error) {
      handleAsyncError(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUsers(0, pageSize).catch(handleAsyncError);
    // fetchGroups().then();
  }, []);

  const manageUser = async (userId, action, record) => {
    try {
      const res = await API.post('/api/user/manage', {
        id: userId,
        action,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess('操作成功完成！');
        let user = res.data.data;
        setUsers(prevUsers => {
          return prevUsers.map(u => {
            if (u.id === userId) {
              if (action === 'delete') {
                return { ...u, status: 'deleted' };
              } else {
                return { ...u, status: user.status, role: user.role };
              }
            }
            return u;
          });
        });
      } else {
        showError(message);
      }
    } catch (error) {
      handleAsyncError(error);
    }
  };

  const searchUsers = async (
    startIdx,
    pageSize,
    searchTitle,
    searchUsername,
  ) => {
    if (searchUsername === '' && setSearchTitle === '') {
      // if keyword is blank, load files instead.
      await loadUsers(startIdx, pageSize);
      return;
    }
    setSearching(true);
    try {
      const res = await API.get(
        `/api/conversation/list_all?title=${searchTitle}&username=${searchUsername}&page=${startIdx}&page_size=${pageSize}`,
      );
      const { success, message, data } = res.data;
      if (success) {
        const newPageData = data.items;
        setActivePage(data.page);
        setUserCount(data.total);
        setUserFormat(newPageData);
      } else {
        showError(message);
      }
    } catch (error) {
      handleAsyncError(error);
    } finally {
      setSearching(false);
    }
  };

  const handleTitleChange = (value) => {
    setSearchTitle(value.trim());
  };

  const handleUsernameChange=(value)=>{
    setSearchUsername(value.trim())
  }

  const handlePageChange = (page) => {
    setActivePage(page);
    if (searchTitle === '' && searchGroup === '') {
      loadUsers(page, pageSize).catch(handleAsyncError);
    } else {
      searchUsers(page, pageSize, searchTitle, searchGroup).catch(handleAsyncError);
    }
  };

  const closeAddUser = () => {
    setShowAddUser(false);
  };

  const closeEditUser = () => {
    setShowEditUser(false);
    setEditingUser({
      id: undefined,
    });
  };

  const refresh = async () => {
    setActivePage(1);
    if (searchTitle === '') {
      await loadUsers(activePage, pageSize);
    } else {
      await searchUsers(activePage, pageSize, searchTitle, searchGroup);
    }
  };

  // const fetchGroups = async () => {
  //   try {
  //     let res = await API.get(`/api/group/`);
  //     // add 'all' option
  //     // res.data.data.unshift('all');
  //     if (res === undefined) {
  //       return;
  //     }
  //     setGroupOptions(
  //       res.data.data.map((group) => ({
  //         label: group,
  //         value: group,
  //       })),
  //     );
  //   } catch (error) {
  //     showError(error.message);
  //   }
  // };

  const handlePageSizeChange = async (size) => {
    console.log("size:", size);
    localStorage.setItem('page-size', size + '');
    setPageSize(size);
    setActivePage(1);
    loadUsers(activePage, size).catch(handleAsyncError);
  };

  return (
    <>
      <Form
        onSubmit={() => {
          searchUsers(activePage, pageSize, searchTitle, searchUsername);
        }}
        labelPosition='left'
      >
        <div style={{ display: 'flex' }}>
          <Space>
            <Tooltip
              content={t('支持搜索用户的 ID、用户名、显示名称和邮箱地址')}
            >
              <Form.Input
                label={t('搜索会话标题')}
                icon='search'
                field='keyword'
                prefix={<IconSearch />}
                placeholder={t('搜索会话标题')}
                value={searchTitle}
                loading={searching}
                onChange={(value) => handleTitleChange(value)}
              />

              <Form.Input
                label={t('搜索  用户')}
                icon='search'
                field='keyword'
                prefix={<IconSearch />}
                placeholder={t('搜索用户')}
                value={searchTitle}
                loading={searching}
                onChange={(value) => handleUsernameChange(value)}
              />
            </Tooltip>

     
            <Button
              label={t('查询')}
              type='primary'
              htmlType='submit'
              className='btn-margin-right'
            >
              {t('查询')}
            </Button>
           
          </Space>
        </div>
        <Table
          columns={columns}
          dataSource={users}
          pagination={{
            formatPageText: (page) =>
              t('第 {{start}} - {{end}} 条，共 {{total}} 条', {
                start: page.currentStart,
                end: page.currentEnd,
                total: users.length,
              }),
            currentPage: activePage,
            pageSize: pageSize,
            total: userCount,
            pageSizeOpts: [10, 20, 50, 100],
            showSizeChanger: true,
            onPageSizeChange: handlePageSizeChange,
            onPageChange: handlePageChange,
          }}
          loading={loading}
        />
      </Form>
    </>
  );
};

export default TableConversation;