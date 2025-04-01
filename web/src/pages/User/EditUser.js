import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API, isMobile, showError, showSuccess } from '../../helpers';
import { renderQuota, renderQuotaWithPrompt } from '../../helpers/render';
import Title from '@douyinfe/semi-ui/lib/es/typography/title';
import {
  TimePicker,
  Button,
  Divider,
  Input,
  Modal,
  Select,
  SideSheet,
  Space,
  Spin,
  Typography,
} from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

const EditUser = (props) => {
  const userId = props.editingUser.id;
  const [loading, setLoading] = useState(true);
  const [addQuotaModalOpen, setIsModalOpen] = useState(false);
  const [addQuotaLocal, setAddQuotaLocal] = useState('');
  const [inputs, setInputs] = useState({
    username: '',
    display_name: '',
    password: '',
    github_id: '',
    oidc_id: '',
    wechat_id: '',
    email: '',
    quota: 0,
    group: 'default',
    start_time_limit:'',
    end_time_limit:''
  });
  const [groupOptions, setGroupOptions] = useState([]);
  const {
    username,
    display_name,
    password,
    github_id,
    oidc_id,
    wechat_id,
    telegram_id,
    email,
    quota,
    group,
  } = inputs;
  const handleInputChange = (name, value) => {
    setInputs((inputs) => ({ ...inputs, [name]: value}));
  };
  const fetchGroups = async () => {
    try {
      let res = await API.get(`/api/group/`);
      setGroupOptions(
        res.data.data.map((group) => ({
          label: group,
          value: group,
        })),
      );
    } catch (error) {
      showError(error.message);
    }
  };
  const navigate = useNavigate();
  const handleCancel = () => {
    props.handleClose();
  };
  const loadUser = async () => {
    setLoading(true);
    let res = undefined;
    if (userId) {
      res = await API.get(`/api/user/${userId}`);
    } else {
      res = await API.get(`/api/user/self`);
    }
    const { success, message, data } = res.data;
    if (success) {
      data.password = '';
      setInputs(data);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadUser().then();
    if (userId) {
      fetchGroups().then();
    }
  }, [props.editingUser.id]);

  const submit = async () => {
    setLoading(true);
    let res = undefined;
    if (userId) {
      let data = { ...inputs, id: parseInt(userId) };
      if (typeof data.quota === 'string') {
        data.quota = parseInt(data.quota);
      }
      res = await API.put(`/api/user/`, data);
    } else {
      res = await API.put(`/api/user/self`, inputs);
    }
    const { success, message } = res.data;
    if (success) {
      showSuccess('用户信息更新成功！');
      props.refresh();
      props.handleClose();
    } else {
      showError(message);
    }
    setLoading(false);
  };
  const onChangeStartTime = (date, dateString) => {
    console.log(date, dateString);
    // setInputs((inputs) => ({ ...inputs, start_time_limit: dateString[0],end_time_limit:end_time_limit }));
    setInputs((inputs) => ({ ...inputs, ['start_time_limit']:timeFormat(`2021-03-27 ${dateString[0]}:00`),['end_time_limit']:timeFormat(`2021-03-27 ${dateString[1]}:00`) }));

    // date 是包含开始时间和结束时间的数组，例如 [startDate, endDate]
    // dateString 是包含开始时间和结束时间的字符串数组，例如 ["2024-01-01", "2024-01-31"]
    
  };

  const donbaToTime=(isoDateStr)=>{
    const date = new Date(isoDateStr);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${hours}:${minutes}`;   
  }

  const timeFormat=(dateStr)=>{
    const date = new Date(dateStr);
// 获取东八区的偏移量（分钟）
    const offset = 8 * 60; 
    // 设置时区偏移量
    date.setMinutes(date.getMinutes() + offset); 
    return date.toISOString().replace('Z', `+08:00`);
   
  }
  const addLocalQuota = () => {
    let newQuota = parseInt(quota) + parseInt(addQuotaLocal);
    setInputs((inputs) => ({ ...inputs, quota: newQuota }));
  };

  const openAddQuotaModal = () => {
    setAddQuotaLocal('0');
    setIsModalOpen(true);
  };

  const { t } = useTranslation();

  return (
    <>
      <SideSheet
        placement={'right'}
        title={<Title level={3}>{t('编辑用户')}</Title>}
        headerStyle={{ borderBottom: '1px solid var(--semi-color-border)' }}
        bodyStyle={{ borderBottom: '1px solid var(--semi-color-border)' }}
        visible={props.visible}
        footer={
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Space>
              <Button theme='solid' size={'large'} onClick={submit}>
                {t('提交')}
              </Button>
              <Button
                theme='solid'
                size={'large'}
                type={'tertiary'}
                onClick={handleCancel}
              >
                {t('取消')}
              </Button>
            </Space>
          </div>
        }
        closeIcon={null}
        onCancel={() => handleCancel()}
        width={isMobile() ? '100%' : 600}
      >
        <Spin spinning={loading}>
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('用户名')}</Typography.Text>
          </div>
          <Input
            label={t('用户名')}
            name='username'
            placeholder={t('请输入新的用户名')}
            onChange={(value) => handleInputChange('username', value)}
            value={username}
            autoComplete='new-password'
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('密码')}</Typography.Text>
          </div>
          <Input
            label={t('密码')}
            name='password'
            mode={'password'}
            placeholder={t('请输入新的密码，最短 8 位')}
            onChange={(value) => handleInputChange('password', value)}
            value={password}
            autoComplete='new-password'
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('显示名称')}</Typography.Text>
          </div>
          <Input
            label={t('显示名称')}
            name='display_name'
            placeholder={t('请输入新的显示名称')}
            onChange={(value) => handleInputChange('display_name', value)}
            value={display_name}
            autoComplete='new-password'
          />
          {userId && (
            <>
              <div style={{ marginTop: 20 }}>
                <Typography.Text>{t('分组')}</Typography.Text>
              </div>
              <Select
                placeholder={t('请选择分组')}
                name='group'
                fluid
                search
                selection
                allowAdditions
                additionLabel={t(
                  '请在系统设置页面编辑分组倍率以添加新的分组：',
                )}
                onChange={(value) => handleInputChange('group', value)}
                value={inputs.group}
                autoComplete='new-password'
                optionList={groupOptions}
              />
              <div style={{ marginTop: 20 }}>
                <TimePicker
                  type='timeRange'
                  format='HH:mm'
                  value={[donbaToTime(inputs['start_time_limit']), donbaToTime(inputs['end_time_limit'])]}
                  onChange={onChangeStartTime}
                />
              </div>
              <div style={{ marginTop: 20 }}>
                <Typography.Text>{`${t('剩余额度')}${renderQuotaWithPrompt(quota)}`}</Typography.Text>
              </div>
              <Space>
                <Input
                  name='quota'
                  placeholder={t('请输入新的剩余额度')}
                  onChange={(value) => handleInputChange('quota', value)}
                  value={quota}
                  type={'number'}
                  autoComplete='new-password'
                />
                <Button onClick={openAddQuotaModal}>{t('添加额度')}</Button>
              </Space>
            </>
          )}
          {/* <Divider style={{ marginTop: 20 }}>{t('以下信息不可修改')}</Divider>
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('已绑定的 GitHub 账户')}</Typography.Text>
          </div>
          <Input
            name='github_id'
            value={github_id}
            autoComplete='new-password'
            placeholder={t('此项只读，需要用户通过个人设置页面的相关绑定按钮进行绑定，不可直接修改')}
            readonly
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('`已绑定的 OIDC 账户')}</Typography.Text>
          </div>
          <Input
              name='oidc_id'
              value={oidc_id}
              placeholder={t('此项只读，需要用户通过个人设置页面的相关绑定按钮进行绑定，不可直接修改')}
              readonly
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('已绑定的微信账户')}</Typography.Text>
          </div>
          <Input
            name='wechat_id'
            value={wechat_id}
            autoComplete='new-password'
            placeholder={t('此项只读，需要用户通过个人设置页面的相关绑定按钮进行绑定，不可直接修改')}
            readonly
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('已绑定的邮箱账户')}</Typography.Text>
          </div>
          <Input
            name='email'
            value={email}
            autoComplete='new-password'
            placeholder={t('此项只读，需要用户通过个人设置页面的相关绑定按钮进行绑定，不可直接修改')}
            readonly
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text>{t('已绑定的Telegram账户')}</Typography.Text>
          </div>
          <Input
            name='telegram_id'
            value={telegram_id}
            autoComplete='new-password'
            placeholder={t('此项只读，需要用户通过个人设置页面的相关绑定按钮进行绑定，不可直接修改')}
            readonly
          /> */}
        </Spin>
      </SideSheet>
      <Modal
        centered={true}
        visible={addQuotaModalOpen}
        onOk={() => {
          addLocalQuota();
          setIsModalOpen(false);
        }}
        onCancel={() => setIsModalOpen(false)}
        closable={null}
      >
        <div style={{ marginTop: 20 }}>
          <Typography.Text>{`${t('新额度')}${renderQuota(quota)} + ${renderQuota(addQuotaLocal)} = ${renderQuota(quota + parseInt(addQuotaLocal))}`}</Typography.Text>
        </div>
        <Input
          name='addQuotaLocal'
          placeholder={t('需要添加的额度（支持负数）')}
          onChange={(value) => {
            setAddQuotaLocal(value);
          }}
          value={addQuotaLocal}
          type={'number'}
          autoComplete='new-password'
        />
      </Modal>
    </>
  );
};

export default EditUser;
