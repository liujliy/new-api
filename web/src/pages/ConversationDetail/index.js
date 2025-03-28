import React, { useContext, useEffect, useState, useCallback } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { API, showError, showSuccess } from '../../helpers';
import { Chat, Button } from '@douyinfe/semi-ui';
import {
  IconArrowLeft,
  IconSidebar,
  IconChevronDown,
} from '@douyinfe/semi-icons';
import './index.css'
const roleInfo = {
  user: {
    name: 'User',
    avatar:
      'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/docs-icon.png',
  },
  assistant: {
    name: 'Assistant',
    avatar:
      'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/other/logo.png',
  },
  system: {
    name: 'System',
    avatar:
      'https://lf3-static.bytednsdoc.com/obj/eden-cn/ptlz_zlp/ljhwZthlaukjlkulzlp/other/logo.png',
  },
};

const ConversationDetail = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { id } = location.state || {};
  const [message, setMessage] = useState();

  useEffect(() => {
    loadConversation();
  }, []);

  


  const detailStr = (str) => {
    const match = str.match(/<think>([\s\S]*?)<\/think>([\s\S]*)/);
    if (match) {
        const contentInside = match[1];
        const contentOutside = match[2];
        const formattedReasoningContent = contentInside.split('\n').map(line => `> ${line}`).join('\n');
        return formattedReasoningContent+contentOutside

    } else {
        return str;
    }
  };

  const loadConversation = async () => {
    const res = await API.get(`/api/conversation/${id}/messages`);
    const { success, message, data } = res.data;
    if (success) {
        setMessage(data);

        setMessage((messages) => {
            return messages.map((element) => {
                try {
                    const parsed = JSON.parse(element.content);
                    // 这里可以对parsed对象进行修改，例如添加一个新属性
                    // 将修改后的对象重新转换为字符串
                    element.content = parsed
                } catch (error) {
                    // 如果解析失败，调用detailStr函数，这里假设detailStr函数已经定义
                    if(typeof(element.content) == "string"){
                        element.content = detailStr(element.content);
                    }
                }
                return element;
            });
        });
        
    } else {
      showError(message);
    }
    // setLoading(false);
  };


  const commonOuterStyle = {
    border: '1px solid var(--semi-color-border)',
    borderRadius: '16px',
    margin: '8px 16px',
    height: 550,
  };
  //输入框为空
  const renderInputArea = useCallback((props) => {
    return <></>;
  }, []);

  return (
    <>
      <Button
        icon={<IconArrowLeft />}
        theme='solid'
        style={{ marginRight: 10 }}
        onClick={() => {
          navigate('/conversation');
        }}
      >
        返回
      </Button>

      <Chat
        mode={'userBubble'}
        renderInputArea={renderInputArea}
        style={commonOuterStyle}
        chats={message}
        roleConfig={roleInfo}
      />
    </>
  );
};

export default ConversationDetail;
